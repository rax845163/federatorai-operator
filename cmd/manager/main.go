package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	fedOperator "github.com/containers-ai/federatorai-operator"
	assets "github.com/containers-ai/federatorai-operator/assets"
	"github.com/containers-ai/federatorai-operator/pkg/apis"
	assetsBin "github.com/containers-ai/federatorai-operator/pkg/assets"
	"github.com/containers-ai/federatorai-operator/pkg/controller"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceread"
	fedOperatorLog "github.com/containers-ai/federatorai-operator/pkg/log"
	alamedaserviceparamter "github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/protocol/grpc"
	"github.com/containers-ai/federatorai-operator/pkg/version"

	certmanagerv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"

	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8sapierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	rest "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const (
	envVarPrefix  = "FEDERATORAI_OPERATOR"
	allowEmptyEnv = true

	defaultLogOutputPath = "/var/log/alameda/federatorai-operator.log"
)

var (
	metricsPort           int32
	configurationFilePath string

	federatoraiOperatorFlagSet = pflag.NewFlagSet("federatorai-operator", pflag.ExitOnError)

	fedOperatorConfig fedOperator.Config

	log = logf.Log.WithName("manager")

	watchNamespace = ""

	registerdAPIResources = make(map[string]bool)
)

func init() {

	initFlags()
	initConfiguration()
	initLogger()
}

func initFlags() {

	federatoraiOperatorFlagSet.Int32Var(&metricsPort, "metrics.port", 8383, "port to export metrics data")
	federatoraiOperatorFlagSet.StringVar(&configurationFilePath, "config", "/etc/federatorai/operator/operator.yml", "File path to federatorai-operator coniguration")

	pflag.CommandLine.AddFlagSet(federatoraiOperatorFlagSet)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()
}

func initConfiguration() {

	fedOperatorConfig = fedOperator.NewDefaultConfig()

	initViperSetting()
	mergeViperValueWithDefaultConfig()
}

func initViperSetting() {

	viper.SetEnvPrefix(envVarPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AllowEmptyEnv(allowEmptyEnv)
	if err := viper.BindPFlags(federatoraiOperatorFlagSet); err != nil {
		panic(err)
	}
}

func mergeViperValueWithDefaultConfig() {

	viper.SetConfigFile(configurationFilePath)

	if err := viper.ReadInConfig(); err != nil {
		panic(errors.New("Read configuration file failed: " + err.Error()))
	}

	if err := viper.Unmarshal(&fedOperatorConfig); err != nil {
		panic(errors.New("Unmarshal configuration failed: " + err.Error()))
	}
}

func initLogger() {

	fedOperatorConfig.Log.AppendOutput(defaultLogOutputPath)
	logger, err := fedOperatorLog.NewZaprLogger(fedOperatorConfig.Log)
	if err != nil {
		panic(err)
	}
	logf.SetLogger(logger)

	fedOperatorConfig.GRPC.Log.AppendOutput(defaultLogOutputPath)
	grpcLogger, err := fedOperatorLog.NewZapLogger(fedOperatorConfig.GRPC.Log)
	if err != nil {
		panic(err)
	}
	grpc.SetGRPCLogger(grpcLogger)
}

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Federatorai Operator Version: %v", version.String))
}

func printConfiguration() {
	if b, err := json.MarshalIndent(fedOperatorConfig, "", "    "); err != nil {
		panic(err.Error())
	} else {
		log.Info(fmt.Sprintf("%+v", string(b)))
	}
}

func main() {

	printVersion()
	printConfiguration()

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	ctx := context.TODO()

	// Become the leader before proceeding
	err = leader.Become(ctx, "federatorai-operator-lock")
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	namespace, found := os.LookupEnv(k8sutil.WatchNamespaceEnvVar)
	if !found {
		namespace = ""
	}

	if err := createCustomeResourceDefinitions(cfg); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := waitCRDReady(cfg); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup requirements before starts the manager
	if err := setupRequirements(cfg); err != nil {
		log.Error(err, "setup requirements failed")
		os.Exit(1)
	}

	//var day time.Duration = 1*24 * time.Hour
	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MetricsBindAddress: fmt.Sprintf("%s:%d", fedOperatorConfig.Metrics.Host, fedOperatorConfig.Metrics.Port),
		//SyncPeriod:         &day,
	})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	if err := securityv1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	if err := autoscaling_v1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	if err := certmanagerv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	if err := apiregistrationv1beta1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Create Service object to expose the metrics port.
	_, err = metrics.ExposeMetricsPort(ctx, metricsPort)
	if err != nil {
		log.Info(err.Error())
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}

func setupRequirements(clientConfig *rest.Config) error {

	if err := createConfigMaps(clientConfig); err != nil {
		return errors.Wrapf(err, "create configMaps failed")
	}

	return nil
}

// In order to let manager has the scheme definition of the crds,
// it should install those crds in to the cluster before the instance of manager is created,
func createCustomeResourceDefinitions(clientConfig *rest.Config) error {

	apiExtensionsClientset, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return errors.Errorf("create k8s clientset failed: %s", err.Error())
	}

	assets := alamedaserviceparamter.GetCustomResourceDefinitions()
	for _, asset := range assets {
		assetBytes, err := assetsBin.Asset(asset)
		if err != nil {
			return errors.Errorf("get asset binary data failed: %s", err.Error())
		}

		crd := resourceread.ReadCustomResourceDefinitionV1Beta1(assetBytes)
		addCRDToRegisterdAPIResources(crd)
		_, err = apiExtensionsClientset.Apiextensions().CustomResourceDefinitions().Create(crd)
		if err != nil && k8sapierrors.IsAlreadyExists(err) {
			log.V(-1).Info("CustomResourceDefinition is existing in cluster, will not create or update it.", "CustomResourceDefinition name", crd.Name)
			continue
		} else if err != nil {
			return errors.Errorf("create CustomResourceDefinition (%s) failed: %s", crd.Name, err.Error())
		}
	}
	return nil
}

func addCRDToRegisterdAPIResources(crd *apiextensionv1beta1.CustomResourceDefinition) {

	group := crd.Spec.Group
	kind := crd.Spec.Names.Kind

	if crd.Spec.Version != "" {
		version := crd.Spec.Version
		groupVersionKind := fmt.Sprintf("%s/%s/%s", group, version, kind)
		addAPIResource(groupVersionKind, registerdAPIResources)
	}

	for _, crdVersion := range crd.Spec.Versions {
		version := crdVersion.Name
		groupVersionKind := fmt.Sprintf("%s/%s/%s", group, version, kind)
		addAPIResource(groupVersionKind, registerdAPIResources)
	}
}

func addAPIResource(groupVersionKind string, gvkMap map[string]bool) {
	gvkMap[groupVersionKind] = true
}

func deleteAPIResource(groupVersionKind string, gvkMap map[string]bool) {
	delete(gvkMap, groupVersionKind)
}

func createConfigMaps(clientConfig *rest.Config) error {

	cli, err := client.New(clientConfig, client.Options{})
	if err != nil {
		return errors.Errorf("new k8s client failed: %s", err.Error())
	}

	ctx := context.TODO()
	files := assets.GetRequiredConfigMaps()
	for _, file := range files {
		fileBytes, err := assetsBin.Asset(file)
		if err != nil {
			return errors.Errorf("get asset's bytes failed: %s", err.Error())
		}
		configMap := resourceread.ReadConfigMapV1(fileBytes)
		err = cli.Create(ctx, configMap)
		if err != nil && !k8sapierrors.IsAlreadyExists(err) {
			return errors.Errorf("create configMap %s/%s failed: %s", configMap.Namespace, configMap.Name, err.Error())
		}
	}

	return nil
}

func waitCRDReady(clientConfig *rest.Config) error {

	waitInterval := 500 * time.Millisecond
	if err := wait.Poll(waitInterval, 30*time.Second, func() (bool, error) {
		apiExtensionsClientset, err := clientset.NewForConfig(clientConfig)
		if err != nil {
			log.V(-1).Info("Create k8s clientset failed, will retry", "msg", err.Error())
			return false, nil
		}

		apiList, err := apiExtensionsClientset.DiscoveryClient.ServerResources()
		if err != nil {
			log.V(-1).Info("Get k8s ServerResources failed, will retry", "msg", err.Error())
			return false, nil
		}

		for _, apiResourceList := range apiList {
			for _, apiResource := range apiResourceList.APIResources {
				groupVersion := apiResourceList.GroupVersion // fmt.Sprintf("%s/%s",group,version)
				kind := apiResource.Kind
				groupVersionKind := fmt.Sprintf("%s/%s", groupVersion, kind)
				deleteAPIResource(groupVersionKind, registerdAPIResources)
			}
		}

		ok := len(registerdAPIResources) == 0
		if !ok {
			log.V(-1).Info("Server does not have required apiResources, will retry fetching")
		}
		return ok, nil

	}); err != nil {
		return err
	}

	return nil
}
