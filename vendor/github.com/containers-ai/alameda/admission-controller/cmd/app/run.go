package app

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	admission_controller "github.com/containers-ai/alameda/admission-controller"
	admission_controller_kubernetes "github.com/containers-ai/alameda/admission-controller/pkg/kubernetes"
	datahub_resource_recommendator "github.com/containers-ai/alameda/admission-controller/pkg/recommendator/resource/datahub"
	"github.com/containers-ai/alameda/admission-controller/pkg/server"
	datahub_controller_validator "github.com/containers-ai/alameda/admission-controller/pkg/validator/controller/datahub"
	utils "github.com/containers-ai/alameda/pkg/utils"
	k8s_utils "github.com/containers-ai/alameda/pkg/utils/kubernetes"
	"github.com/containers-ai/alameda/pkg/utils/kubernetes/metadata"
	"github.com/containers-ai/alameda/pkg/utils/log"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"

	openshiftappsv1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	admissionregistration_v1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	sigs_k8s_client "sigs.k8s.io/controller-runtime/pkg/client"
	sigs_k8s_client_config "sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	envVarPrefix  = "ALAMEDA_ADMCTL"
	allowEmptyEnv = true

	defaultRotationMaxSizeMegabytes = 100
	defaultRotationMaxBackups       = 7
	defaultLogRotateOutputFile      = "/var/log/alameda/alameda-adm-ctrl.log"
)

var (
	scope         = log.RegisterScope("admission-controller run", "admission-controller", 0)
	config        *admission_controller.Config
	k8sRestConfig *rest.Config

	mutatingWebhookConfigurationInstance admissionregistration_v1beta1.MutatingWebhookConfiguration

	ownerReferenceOfControllerOwningAdmissionController meta_v1.OwnerReference

	sigsK8SClient        sigs_k8s_client.Client
	datahubServiceClient datahub_v1alpha1.DatahubServiceClient

	RunCmd = &cobra.Command{
		Use:              "run",
		Short:            "start alameda admission-controller server",
		Long:             "",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			flag.Parse()

			initConfig()
			initLog()
			initThirdPartyClient()

			if !config.Enable {
				scope.Warn("admission-controller is not enabled, os exit with 0")
				os.Exit(0)
			}

			if err := initOwnerReferenceOfControllerOwningAdmissionController(); err != nil {
				panic(err)
			}

			if err := prepareRequirements(); err != nil {
				panic(err)
			}
			if err := setupRequirements(); err != nil {
				panic(err)
			}

			datahubResourceRecommendator, err := datahub_resource_recommendator.NewDatahubResourceRecommendator(datahubServiceClient)
			if err != nil {
				panic(err.Error())
			}
			datahubControllerValidator := datahub_controller_validator.NewControllerValidator(datahubServiceClient, sigsK8SClient)
			admissionController, err := server.NewAdmissionControllerWithConfig(server.Config{Enable: config.Enable}, sigsK8SClient, datahubResourceRecommendator, datahubControllerValidator)
			if err != nil {
				panic(err.Error())
			}

			mux := http.NewServeMux()
			registerHandlerFunc(mux, admissionController)

			server := newHTTPServer(*config, mux)
			if err := server.ListenAndServeTLS("", ""); err != nil {
				panic(err.Error())
			}

			return nil
		},
	}
)

func initConfig() {

	defaultConfig := admission_controller.NewDefaultConfig()
	config = &defaultConfig
	initViperSetting()
	mergeConfigFileValueWithDefaultConfigValue()

	k8sCfg, err := sigs_k8s_client_config.GetConfig()
	if err != nil {
		panic(errors.Errorf("initialization failed: get k8s rest configuration failed: %s", err.Error()))
	}
	k8sRestConfig = k8sCfg

}

func initViperSetting() {

	viper.SetEnvPrefix(envVarPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// viper.AllowEmptyEnv(allowEmptyEnv)
}

func mergeConfigFileValueWithDefaultConfigValue() {

	viper.SetConfigFile(configurationFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(errors.New("Read configuration file failed: " + err.Error()))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(errors.New("Unmarshal configuration failed: " + err.Error()))
	}
}

func initLog() {

	opt := log.DefaultOptions()
	opt.RotationMaxSize = defaultRotationMaxSizeMegabytes
	opt.RotationMaxBackups = defaultRotationMaxBackups
	opt.RotateOutputPath = defaultLogRotateOutputFile
	err := log.Configure(opt)
	if err != nil {
		panic(err)
	}

	for _, scope := range log.Scopes() {
		scope.SetLogCallers(config.Log.SetLogCallers == true)
		if outputLvl, ok := log.StringToLevel(config.Log.OutputLevel); ok {
			scope.SetOutputLevel(outputLvl)
		}
		if stacktraceLevel, ok := log.StringToLevel(config.Log.StackTraceLevel); ok {
			scope.SetStackTraceLevel(stacktraceLevel)
		}
	}
}

func initThirdPartyClient() error {

	var err error

	conn, err := grpc.Dial(config.Datahub.Address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	datahubServiceClient = datahub_v1alpha1.NewDatahubServiceClient(conn)

	sigsK8SClient, err = sigs_k8s_client.New(k8sRestConfig, sigs_k8s_client.Options{Scheme: admission_controller_kubernetes.Scheme})
	if err != nil {
		return err
	}

	return nil
}

func initOwnerReferenceOfControllerOwningAdmissionController() error {

	pod, err := k8s_utils.GetPodByNamespaceNameWithConfig(utils.GetRunningNamespace(), utils.GetRunnningPodName(), *k8sRestConfig)
	if err != nil {
		return errors.Wrap(err, "initialization failed: get deployment or deploymentconfig owning admission-controller failed")
	}

	ort, err := metadata.NewOwnerReferenceTracerWithConfig(*k8sRestConfig)
	if err != nil {
		return errors.Wrap(err, "initialization failed: get deployment or deploymentconfig owning admission-controller failed")
	}

	dep, dc, err := ort.GetDeploymentOrDeploymentConfigOwningPod(pod)
	if err != nil {
		return errors.Wrapf(err, "initialization failed: get deployment or deploymentconfig owning admission-controller failed")
	}

	if dep == nil && dc == nil {
		return errors.New("initialization failed: cannot get deployment or deploymentConfig owning admission-controller")
	}

	var ownerType meta_v1.TypeMeta
	var ownerMeta meta_v1.ObjectMeta
	if dep != nil {
		ownerType.APIVersion = appsv1.SchemeGroupVersion.String()
		ownerType.Kind = "Deployment"
		ownerMeta = dep.ObjectMeta
	} else if dc != nil {
		ownerType.APIVersion = openshiftappsv1.SchemeGroupVersion.String()
		ownerType.Kind = "DeploymentConfig"
		ownerMeta = dc.ObjectMeta
	}
	ownerReferenceOfControllerOwningAdmissionController = k8s_utils.NewOwnerReference(ownerType, ownerMeta, false)

	return nil
}

func prepareRequirements() error {

	scope.Debugf("preparing requirements")

	err := prepareMutatingWebhookConfigurationInstance()
	if err != nil {
		return errors.Wrap(err, "prepare requirements failed")
	}
	return nil
}

func prepareMutatingWebhookConfigurationInstance() error {

	scope.Debugf("preparing MutatingWebhookConfiguration instance")

	var (
		namespace         string
		caBundle          []byte
		serviceName       = config.Service.Name
		mutatePodEndPoint = "/pods"
	)

	namespace = config.DeployedNamespace

	caBundle, err := ioutil.ReadFile(config.CACertFile)
	if err != nil {
		return errors.Errorf("prepare MutatingWebhookConfiguration failed: read caBundle file failed: %s", err.Error())
	}

	mutatingWebhookConfigurationInstance = admissionregistration_v1beta1.MutatingWebhookConfiguration{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "MutatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1beta1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:            fmt.Sprintf("mutating-webhook.admission-controller.%s.svc", namespace),
			OwnerReferences: []meta_v1.OwnerReference{ownerReferenceOfControllerOwningAdmissionController},
		},
		Webhooks: []admissionregistration_v1beta1.Webhook{
			admissionregistration_v1beta1.Webhook{
				Name: fmt.Sprintf("mutating-webhook.admission-controller.%s.svc", namespace),
				Rules: []admissionregistration_v1beta1.RuleWithOperations{
					admissionregistration_v1beta1.RuleWithOperations{
						Operations: []admissionregistration_v1beta1.OperationType{
							admissionregistration_v1beta1.Create,
						},
						Rule: admissionregistration_v1beta1.Rule{
							APIGroups: []string{
								"",
							},
							APIVersions: []string{
								"v1",
							},
							Resources: []string{
								"pods",
							},
						},
					},
				},
				ClientConfig: admissionregistration_v1beta1.WebhookClientConfig{
					CABundle: caBundle,
					Service: &admissionregistration_v1beta1.ServiceReference{
						Namespace: namespace,
						Name:      serviceName,
						Path:      &mutatePodEndPoint,
					},
				},
			},
		},
	}

	return nil
}

func setupRequirements() error {

	scope.Debugf("setting up requirements")

	err := createOrUpdateMutatingWebhookConfiguration(mutatingWebhookConfigurationInstance)
	if err != nil {
		return errors.Wrap(err, "setup requirements failed")
	}

	return nil
}

func createOrUpdateMutatingWebhookConfiguration(instance admissionregistration_v1beta1.MutatingWebhookConfiguration) error {

	scope.Debugf("setting up MutatingWebhookConfiguration")

	ctx := context.TODO()
	currentInstance := admissionregistration_v1beta1.MutatingWebhookConfiguration{}
	err := sigsK8SClient.Get(ctx, sigs_k8s_client.ObjectKey{Name: instance.Name}, &currentInstance)
	if err != nil && k8s_errors.IsNotFound(err) {

		scope.Debugf("no existing MutatingWebhookConfiguration: %s, create new one", instance.Name)

		err = sigsK8SClient.Create(ctx, &instance)
		if err != nil {
			return errors.Errorf("create MutatingWebhookConfigurations failed: %s", err.Error())
		}
	} else if err != nil {

		return errors.Errorf("get MutatingWebhookConfigurations failed: %s", err.Error())
	} else {

		scope.Debugf("found existing MutatingWebhookConfiguration, update. Previous: %+v, Updated: %+v .", currentInstance, instance)

		currentInstance.Webhooks = instance.Webhooks
		currentInstance.ObjectMeta.OwnerReferences = []meta_v1.OwnerReference{ownerReferenceOfControllerOwningAdmissionController}
		err = sigsK8SClient.Update(ctx, &currentInstance)
		if err != nil {
			return errors.Errorf("update MutatingWebhookConfigurations failed: %s", err.Error())
		}
	}

	return nil
}

func registerHandlerFunc(mux *http.ServeMux, ac server.AdmissionController) {
	mux.HandleFunc("/pods", ac.MutatePod)
}

func newHTTPServer(cfg admission_controller.Config, mux *http.ServeMux) *http.Server {
	tlsConfig, err := config.ConfigTLS()
	if err != nil {
		panic(err.Error())
	}

	addr := fmt.Sprintf(":%s", fmt.Sprint(config.Port))
	server := &http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	return server
}
