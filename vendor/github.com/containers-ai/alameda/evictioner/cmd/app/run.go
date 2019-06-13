package app

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/containers-ai/alameda/cmd/app"
	"github.com/containers-ai/alameda/evictioner/pkg/eviction"
	"github.com/containers-ai/alameda/operator/pkg/apis"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8s_config "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	RunCmd = &cobra.Command{
		Use:   "run",
		Short: "start alameda evictioner",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			app.PrintSoftwareVer()
			initConfig()
			initLogger()
			setLoggerScopesWithConfig(*config.Log)
			displayConfig()
			startEvictioner()
		},
	}
)

func displayConfig() {
	if configBin, err := json.MarshalIndent(config, "", "  "); err != nil {
		scope.Error(err.Error())
	} else {
		scope.Infof(fmt.Sprintf("Evict configuration: %s", string(configBin)))
	}
}

func startEvictioner() {
	conn, err := grpc.Dial(config.Datahub.Address, grpc.WithInsecure())
	if err != nil {
		scope.Errorf("create pods to datahub failed: %s", err.Error())
		return
	}

	defer conn.Close()

	datahubServiceClnt := datahub_v1alpha1.NewDatahubServiceClient(conn)

	k8sClientConfig, err := k8s_config.GetConfig()
	if err != nil {
		scope.Error("Get kubernetes configuration failed: " + err.Error())
		return
	}

	k8sCli, err := client.New(k8sClientConfig, client.Options{})
	if err != nil {
		scope.Error("Create kubernetes client failed: " + err.Error())
		return
	}

	mgr, err := manager.New(k8sClientConfig, manager.Options{})
	if err != nil {
		scope.Error(err.Error())
	}
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		scope.Error(err.Error())
	}

	evictioner := eviction.NewEvictioner(config.Eviction.CheckCycle,
		datahubServiceClnt,
		k8sCli,
		*config.Eviction,
	)
	evictioner.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
