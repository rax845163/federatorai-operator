package component

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"html/template"
	"net"
	"strconv"
	"strings"

	autoscaling_v1alpha1 "github.com/containers-ai/alameda/operator/pkg/apis/autoscaling/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/assets"
	"github.com/containers-ai/federatorai-operator/pkg/lib/resourceread"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ingressv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/client-go/util/cert"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_alamedaservice")

// assetsSkipTemplateParse contains asset files that cannnot do template parsing
var assetsSkipTemplateParse = []string{
	alamedaserviceparamter.ConfigMapDashboardsConfig,
}

type ComponentConfig struct {
	NameSpace         string
	PodTemplateConfig PodTemplateConfig
}

func NewComponentConfig(namespace string, ptc PodTemplateConfig) *ComponentConfig {

	return &ComponentConfig{
		NameSpace:         namespace,
		PodTemplateConfig: ptc,
	}
}

func (c *ComponentConfig) SetNameSpace(ns string) {
	c.NameSpace = ns
}

func (c ComponentConfig) templateAssets(data string) []byte {
	tmpl, err := template.New("namespaceServiceToYaml").Parse(data)
	if err != nil {
		panic(err)
	}
	yamlBuffer := new(bytes.Buffer)
	if err = tmpl.Execute(yamlBuffer, c); err != nil {
		panic(err)
	}
	return yamlBuffer.Bytes()
}

func (c ComponentConfig) NewClusterRoleBinding(str string) *rbacv1.ClusterRoleBinding {
	crbByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create clusterrolebinding")

	}
	crb := resourceread.ReadClusterRoleBindingV1(c.templateAssets(string(crbByte[:])))
	return crb
}
func (c ComponentConfig) NewClusterRole(str string) *rbacv1.ClusterRole {
	crByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create clusterrole")
	}
	cr := resourceread.ReadClusterRoleV1(c.templateAssets(string(crByte[:])))
	return cr
}
func (c ComponentConfig) NewServiceAccount(str string) *corev1.ServiceAccount {
	saByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create serviceaccount")

	}
	sa := resourceread.ReadServiceAccountV1(c.templateAssets(string(saByte[:])))
	return sa
}
func (c ComponentConfig) NewConfigMap(str string) *corev1.ConfigMap {
	cmByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create configmap")
	}

	var cm *corev1.ConfigMap
	if skipTemplateParse(str) {
		cm = resourceread.ReadConfigMapV1(cmByte)
		cm.Namespace = c.NameSpace
	} else {
		cm = resourceread.ReadConfigMapV1(c.templateAssets(string(cmByte[:])))
	}
	return cm
}
func (c ComponentConfig) NewPersistentVolumeClaim(str string) *corev1.PersistentVolumeClaim {
	pvcByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create persistentvolumeclaim")

	}
	pvc := resourceread.ReadPersistentVolumeClaimV1(c.templateAssets(string(pvcByte[:])))
	return pvc
}
func (c ComponentConfig) NewService(str string) *corev1.Service {
	svByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create service")

	}
	sv := resourceread.ReadServiceV1(c.templateAssets(string(svByte[:])))
	return sv
}

func (c ComponentConfig) NewAlamedaScaler(str string) *autoscaling_v1alpha1.AlamedaScaler {
	scalerByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create NewAlamedaScaler")

	}
	scaler := resourceread.ReadScalerV1(c.templateAssets(string(scalerByte[:])))
	return scaler
}

func (c ComponentConfig) NewDeployment(str string) *appsv1.Deployment {
	deploymentBytes, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create deployment")

	}
	d := resourceread.ReadDeploymentV1(c.templateAssets(string(deploymentBytes[:])))
	d.Spec.Template = c.mutatePodTemplateSpecWithConfig(d.Spec.Template)
	return d
}

func (c ComponentConfig) NewRoute(str string) *routev1.Route {
	rtByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create route")

	}
	rt := resourceread.ReadRouteV1(c.templateAssets(string(rtByte[:])))
	return rt
}

func (c ComponentConfig) NewIngress(str string) *ingressv1beta1.Ingress {
	igByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create ingress")

	}
	ig := resourceread.ReadIngressv1beta1(c.templateAssets(string(igByte[:])))
	return ig
}

func (c ComponentConfig) NewStatefulSet(str string) *appsv1.StatefulSet {
	ssByte, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create statefulset")

	}
	ss := resourceread.ReadStatefulSetV1(c.templateAssets(string(ssByte[:])))
	return ss
}

func (c ComponentConfig) NewAdmissionControllerSecret() (*corev1.Secret, error) {

	secret, err := c.NewSecret("Secret/admission-controller-tls.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to buiild admission-controller secret")
	}

	caKey, err := cert.NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "new ca private key failed")
	}

	caCertCfg := cert.Config{}
	caCert, err := cert.NewSelfSignedCACert(caCertCfg, caKey)
	if err != nil {
		return nil, errors.Wrap(err, "new ca cert failed")
	}

	admctlKey, err := cert.NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "new admctl private key failed")
	}

	admctlCertCfg := cert.Config{
		CommonName: fmt.Sprintf("admission-controller.%s.svc", c.NameSpace),
		Usages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
	}
	admctlCert, err := cert.NewSignedCert(admctlCertCfg, admctlKey, caCert, caKey)
	if err != nil {
		return nil, errors.Wrap(err, "new admctl cert failed")
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data["ca.crt"] = cert.EncodeCertPEM(caCert)
	secret.Data["tls.crt"] = cert.EncodeCertPEM(admctlCert)
	secret.Data["tls.key"] = cert.EncodePrivateKeyPEM(admctlKey)

	return secret, nil
}

func (c ComponentConfig) NewfedemeterSecret() (*corev1.Secret, error) {
	secret, err := c.NewSecret("Secret/fedemeter-tls.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to buiild fedemeter secret")
	}
	host := fmt.Sprintf("fedemeter-api.%s.svc", c.NameSpace)
	crt, key, err := cert.GenerateSelfSignedCertKey(host, []net.IP{}, []string{})
	if err != nil {
		return nil, errors.Errorf("failed to buiild fedemeter secret: %s", err.Error())
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data["tls.crt"] = crt
	secret.Data["tls.key"] = key
	return secret, nil
}

func (c ComponentConfig) NewInfluxDBSecret() (*corev1.Secret, error) {

	secret, err := c.NewSecret("Secret/alameda-influxdb.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to buiild influxdb secret")
	}

	host := fmt.Sprintf("admission-influxdb.%s.svc", c.NameSpace)
	crt, key, err := cert.GenerateSelfSignedCertKey(host, []net.IP{}, []string{})
	if err != nil {
		return nil, errors.Errorf("failed to buiild influxdb secret: %s", err.Error())
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data["tls.crt"] = crt
	secret.Data["tls.key"] = key

	return secret, nil
}

func (c ComponentConfig) NewSecret(str string) (*corev1.Secret, error) {
	secretBytes, err := assets.Asset(str)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build secret from assets' bin data")
	}
	s, err := resourceread.ReadSecretV1(c.templateAssets(string(secretBytes[:])))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build secret from assets' bin data")
	}
	return s, nil
}

func (c ComponentConfig) RegistryCustomResourceDefinition(str string) *apiextv1beta1.CustomResourceDefinition {
	crdBytes, err := assets.Asset(str)
	if err != nil {
		log.Error(err, "Failed to Test create testcrd")
	}
	crd := resourceread.ReadCustomResourceDefinitionV1Beta1(crdBytes)
	return crd
}

func (c ComponentConfig) mutatePodTemplateSpecWithConfig(podTemplateSpec corev1.PodTemplateSpec) corev1.PodTemplateSpec {

	copyPodTemplateSpec := podTemplateSpec.DeepCopy()

	var currentPodSecurityContext corev1.PodSecurityContext
	if copyPodTemplateSpec.Spec.SecurityContext != nil {
		currentPodSecurityContext = *copyPodTemplateSpec.Spec.SecurityContext
	}
	podSecurityContext := c.mutatePodSecurityContextWithConfig(currentPodSecurityContext)
	copyPodTemplateSpec.Spec.SecurityContext = &podSecurityContext

	return *copyPodTemplateSpec
}

func (c ComponentConfig) mutatePodSecurityContextWithConfig(podSecurityContext corev1.PodSecurityContext) corev1.PodSecurityContext {

	copyPodSecurityContext := podSecurityContext.DeepCopy()

	if c.PodTemplateConfig.PodSecurityContext.FSGroup != nil {
		fsGroup := *c.PodTemplateConfig.PodSecurityContext.FSGroup
		copyPodSecurityContext.FSGroup = &fsGroup
	}

	return *copyPodSecurityContext
}

// PodTemplateConfig specifies pod confiruation needed while deploying pod
type PodTemplateConfig struct {
	corev1.PodSecurityContext
}

func NewDefaultPodTemplateConfig(ns corev1.Namespace) PodTemplateConfig {

	var (
		ptc PodTemplateConfig

		defaultPSC         corev1.PodSecurityContext
		okdPreAllocatedPSC corev1.PodSecurityContext
	)

	defaultPSC = newDefaultPodSecurityContext()
	ptc = PodTemplateConfig{
		PodSecurityContext: defaultPSC,
	}

	okdPreAllocatedPSC = newOKDPreAllocatedPodSecurityContext(ns)
	ptc.PodSecurityContext = overwritePodSecurityContextFromOKDPodSecurityContext(ptc.PodSecurityContext, okdPreAllocatedPSC)

	return ptc
}

func newDefaultPodSecurityContext() corev1.PodSecurityContext {

	var (
		defaultFSGroup = int64(1000)
	)

	psc := corev1.PodSecurityContext{
		FSGroup: &defaultFSGroup,
	}

	return psc
}

// Currently implement fsGroup strategy.
// Please reference okd documentation https://docs.okd.io/latest/architecture/additional_concepts/authorization.html#understanding-pre-allocated-values-and-security-context-constraints
func newOKDPreAllocatedPodSecurityContext(ns corev1.Namespace) corev1.PodSecurityContext {

	var psc corev1.PodSecurityContext

	annotations := ns.GetObjectMeta().GetAnnotations()

	var fsGroup *int64
	minFSGroupValueString := ""
	if fsGroupRanges, exist := annotations["openshift.io/sa.scc.supplemental-groups"]; exist {
		firstFSGroupRange := strings.Split(fsGroupRanges, ",")[0]
		if strings.Contains(firstFSGroupRange, "/") {
			minFSGroupValueString = strings.Split(firstFSGroupRange, "/")[0]
		} else if strings.Contains(firstFSGroupRange, "-") {
			minFSGroupValueString = strings.Split(firstFSGroupRange, "-")[0]
		}
	} else if fsGroupRange, exist := annotations["openshift.io/sa.scc.uid-range"]; exist {
		if strings.Contains(fsGroupRange, "/") {
			minFSGroupValueString = strings.Split(fsGroupRange, "/")[0]
		}
	}
	if minFSGroupValueString != "" {
		if minFSGroupValue, err := strconv.ParseInt(minFSGroupValueString, 10, 64); err != nil {
			log.V(-1).Info("parse minimum fsGroup value from namespace's annotation failed", "errMsg", err.Error())
		} else {
			fsGroup = &minFSGroupValue
		}
	}
	psc.FSGroup = fsGroup

	return psc
}

// Currently overwrite fsGroup
// Please reference okd documentation https://docs.okd.io/latest/architecture/additional_concepts/authorization.html#understanding-pre-allocated-values-and-security-context-constraints
func overwritePodSecurityContextFromOKDPodSecurityContext(psc, okdPSC corev1.PodSecurityContext) corev1.PodSecurityContext {

	copyPSC := psc.DeepCopy()
	copyOKDPSC := okdPSC.DeepCopy()

	copyPSC.FSGroup = copyOKDPSC.FSGroup

	return *copyPSC
}

func skipTemplateParse(asset string) bool {
	return util.StringInSlice(asset, assetsSkipTemplateParse)
}
