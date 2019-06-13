/*
Copyright 2019 The Alameda Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/containers-ai/alameda/operator"

	"github.com/containers-ai/alameda/operator/pkg/apis"
	"github.com/containers-ai/alameda/operator/pkg/controller"
	"github.com/containers-ai/alameda/operator/pkg/probe"
	"github.com/containers-ai/alameda/operator/pkg/webhook"
	logUtil "github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const (
	envVarPrefix = "ALAMEDA_OPERATOR"

	defaultRotationMaxSizeMegabytes = 100
	defaultRotationMaxBackups       = 7
	defaultLogRotateOutputFile      = "/var/log/alameda/alameda-operator.log"
)

const JSONIndent = "  "

var operatorConfigFile string
var crdLocation string
var showVer bool
var readinessProbeFlag bool
var livenessProbeFlag bool

var operatorConf operator.Config
var k8sConfig *rest.Config
var scope *logUtil.Scope

var (
	// VERSION is sofeware version
	VERSION string
	// BUILD_TIME is build time
	BUILD_TIME string
	// GO_VERSION is go version
	GO_VERSION string
)

func init() {
	flag.BoolVar(&showVer, "version", false, "show version")
	flag.BoolVar(&readinessProbeFlag, "readiness-probe", false, "probe for readiness")
	flag.BoolVar(&livenessProbeFlag, "liveness-probe", false, "probe for liveness")
	flag.StringVar(&operatorConfigFile, "config", "/etc/alameda/operator/operator.yml", "File path to operator coniguration")
	flag.StringVar(&crdLocation, "crd-location", "/etc/alameda/operator/crds", "CRD location")

	scope = logUtil.RegisterScope("manager", "operator entry point", 0)
}

func initLogger() {

	opt := logUtil.DefaultOptions()
	opt.RotationMaxSize = defaultRotationMaxSizeMegabytes
	opt.RotationMaxBackups = defaultRotationMaxBackups
	opt.RotateOutputPath = defaultLogRotateOutputFile
	err := logUtil.Configure(opt)
	if err != nil {
		panic(err)
	}

	scope.Infof("Log output level is %s.", operatorConf.Log.OutputLevel)
	scope.Infof("Log stacktrace level is %s.", operatorConf.Log.StackTraceLevel)
	for _, scope := range logUtil.Scopes() {
		scope.SetLogCallers(operatorConf.Log.SetLogCallers == true)
		if outputLvl, ok := logUtil.StringToLevel(operatorConf.Log.OutputLevel); ok {
			scope.SetOutputLevel(outputLvl)
		}
		if stacktraceLevel, ok := logUtil.StringToLevel(operatorConf.Log.StackTraceLevel); ok {
			scope.SetStackTraceLevel(stacktraceLevel)
		}
	}
}

func initServerConfig(mgr *manager.Manager) {

	operatorConf = operator.NewConfigWithoutMgr()
	if mgr != nil {
		operatorConf = operator.NewConfig(*mgr)
	}

	viper.SetEnvPrefix(envVarPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// TODO: This config need default value. And it should check the file exists befor SetConfigFile.
	viper.SetConfigFile(operatorConfigFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(errors.New("Read configuration failed: " + err.Error()))
	}
	err = viper.Unmarshal(&operatorConf)
	if err != nil {
		panic(errors.New("Unmarshal configuration failed: " + err.Error()))
	} else {
		if operatorConfBin, err := json.MarshalIndent(operatorConf, "", JSONIndent); err == nil {
			scope.Infof(fmt.Sprintf("Operator configuration: %s", string(operatorConfBin)))
		}
	}
}

func main() {
	flag.Parse()
	if showVer {
		printSoftwareInfo()
		return
	}

	if readinessProbeFlag && livenessProbeFlag {
		scope.Error("Cannot run readiness probe and liveness probe at the same time")
		return
	} else if readinessProbeFlag {
		initServerConfig(nil)
		readinessProbe(&probe.ReadinessProbeConfig{
			ValidationSrvPort: operatorConf.K8SWebhookServer.Port,
			DatahubAddr:       operatorConf.Datahub.Address,
		})
		return
	} else if livenessProbeFlag {
		initServerConfig(nil)
		livenessProbe(&probe.LivenessProbeConfig{
			ValidationSvc: &probe.ValidationSvc{
				SvcName: operatorConf.K8SWebhookServer.Service.Name,
				SvcNS:   operatorConf.K8SWebhookServer.Service.Namespace,
				SvcPort: operatorConf.K8SWebhookServer.Service.Port,
			},
		})
		return
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		scope.Error("Get configuration failed: " + err.Error())
	}
	k8sConfig = cfg

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(k8sConfig, manager.Options{})
	if err != nil {
		scope.Error(err.Error())
	}

	// TODO: There are config dependency, this manager should have it's config.
	initServerConfig(&mgr)
	initLogger()
	printSoftwareInfo()

	scope.Info("Registering Components.")
	registerThirdPartyCRD()
	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		scope.Error(err.Error())
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		scope.Error(err.Error())
	}

	scope.Info("setting up webhooks")
	if err := webhook.AddToManager(mgr); err != nil {
		scope.Errorf("unable to register webhooks to the manager: %s", err.Error())
		os.Exit(1)
	}

	// go applyCRDs(k8sConfig)
	go registerNodes(mgr.GetClient(), operatorConf.Datahub.RetryInterval.Default)
	go syncAlamedaPodsWithDatahub(mgr.GetClient(), operatorConf.Datahub.RetryInterval.Default)
	go launchWebhook(&mgr, &operatorConf)
	go addOwnerReferenceToResourcesCreateFrom3rdPkg(mgr.GetClient())
	scope.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		scope.Error(err.Error())
	}
}
