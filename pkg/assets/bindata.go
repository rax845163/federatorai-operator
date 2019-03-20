package assets

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}
type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

var _bindata = map[string]func() (*asset, error){
	"../../assets/CustomResourceDefinition/alamedarecommendationsCRD.yaml": alamedarecommendationsCRDYaml,
	"../../assets/CustomResourceDefinition/alamedascalersCRD.yaml":         alamedascalersCRDYaml,
	"../../assets/ClusterRoleBinding/alameda-datahubCRB.yaml":              alamedadatahubCRBYaml,
	"../../assets/ClusterRoleBinding/alameda-operatorCRB.yaml":             alamedaoperatorCRBYaml,
	"../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml":           alamedaevictionerCRBYaml,
	"../../assets/ClusterRoleBinding/admission-controllerCRB.yaml":         admissioncontrollerCRBYaml,
	"../../assets/ClusterRole/alameda-datahubCR.yaml":                      alamedadatahubCRYaml,
	"../../assets/ClusterRole/alameda-operatorCR.yaml":                     alamedaoperatorCRYaml,
	"../../assets/ClusterRole/alameda-evictionerCR.yaml":                   alamedaevictionerCRYaml,
	"../../assets/ClusterRole/admission-controllerCR.yaml":                 admissioncontrollerCRYaml,
	"../../assets/ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml": aggregatealamedaadmineditalamedaCRYaml,
	"../../assets/ServiceAccount/alameda-datahubSA.yaml":                   alamedadatahubSAYaml,
	"../../assets/ServiceAccount/alameda-operatorSA.yaml":                  alamedaoperatorSAYaml,
	"../../assets/ServiceAccount/alameda-evictionerSA.yaml":                alamedaevictionerSAYaml,
	"../../assets/ServiceAccount/admission-controllerSA.yaml":              admissioncontrollerSAYaml,
	"../../assets/ConfigMap/grafana-datasources.yaml":                      grafanadatasourcesCMYaml,
	"../../assets/PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml":        myalamedainfluxdbPVCYaml,
	"../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml":         myalamedagrafanaPVCYaml,
	"../../assets/Service/alameda-datahubSV.yaml":                          alamedadatahubSVYaml,
	"../../assets/Service/admission-controllerSV.yaml":                     admissioncontrollerSVYaml,
	"../../assets/Service/alameda-influxdbSV.yaml":                         alamedainfluxdbSVYaml,
	"../../assets/Service/alameda-grafanaSV.yaml":                          alamedagrafanaSVYaml,
	"../../assets/Deployment/alameda-datahubDM.yaml":                       alamedadatahubDMYaml,
	"../../assets/Deployment/alameda-operatorDM.yaml":                      alamedaoperatorDMYaml,
	"../../assets/Deployment/alameda-evictionerDM.yaml":                    alamedaevictionerDMYaml,
	"../../assets/Deployment/admission-controllerDM.yaml":                  admissioncontrollerDMYaml,
	"../../assets/Deployment/alameda-influxdbDM.yaml":                      alamedainfluxdbDMYaml,
	"../../assets/Deployment/alameda-grafanaDM.yaml":                       alamedagrafanaDMYaml,
	/**********************************************************************************************/

	"../../manifests/TestCrds.yaml":              manifestsCronTabcrdYaml,
	"../../manifests/deployment.yaml":            manifestsDeploymentYaml,
	"../../manifests/serviceaccount.yaml":        manifestsServiceaccountYaml,
	"../../manifests/configmap.yaml":             manifestsConfigMapYaml,
	"../../manifests/secret.yaml":                manifestsSecretTokenYaml,
	"../../manifests/service.yaml":               manifestsServiceYaml,
	"../../manifests/clusterrole.yaml":           manifestsClusterRoleYaml,
	"../../manifests/clusterrolebinding.yaml":    manifestsClusterRoleBindingYaml,
	"../../manifests/persistentvolumeclaim.yaml": manifestsPersistentVolumeClaimYaml,
}

func ReadYamlFile(file_location string) []byte {
	file, err := os.Open(file_location)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := fileinfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	return buffer
}

var _manifestsCronTabcrdYaml = []byte(`apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
spec:
  group: stable.example.com
  names:
    kind: CronTab
    listKind: CronTabList
    plural: crontabs
    singular: crontabs
  scope: Cluster
  subresources:
    status: {}
  versions:
    - name: v1
      served: true
      storage: true
`)

func manifestsCronTabcrdYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/TestCrds.yaml"), nil
	//return _manifestsCronTabcrdYaml, nil
}
func manifestsCronTabcrdYaml() (*asset, error) {
	bytes, err := manifestsCronTabcrdYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/TestCrds.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _manifestsDeploymentYamlYaml = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: webtest777
  namespace: kroos-tutorial
  labels:
    app: web
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - image: nginx
        name: nginxweb
        ports:
        - containerPort: 8080
          name: http
        resources:
          requests:
            cpu: 10m
            memory: 50Mi`)

func manifestsDeploymentYamlYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/deployment.yaml"), nil
	//return _manifestsDeploymentYamlYaml, nil
}

func manifestsDeploymentYaml() (*asset, error) {
	bytes, err := manifestsDeploymentYamlYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/deployment.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _manifestsServiceaccountYaml = []byte(`apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: kroos-tutorial
  name: kroos-controller
`)

func manifestsServiceaccountYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/serviceaccount.yaml"), nil
	//return _manifestsServiceaccountYaml, nil
}

func manifestsServiceaccountYaml() (*asset, error) {
	bytes, err := manifestsServiceaccountYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/serviceaccount.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _manifestsSecretTokenYaml = []byte(`apiVersion: v1
kind: Secret
metadata:
  annotations:
    kubernetes.io/service-account.name: kroos-controller
  name: kroos-controller-token
  namespace: kroos-tutorial
type: kubernetes.io/service-account-token
`)

func manifestsSecretTokenYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/secret.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}

func manifestsSecretTokenYaml() (*asset, error) {
	bytes, err := manifestsSecretTokenYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/secret.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func manifestsConfigMapYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/configmap.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}

func manifestsConfigMapYaml() (*asset, error) {
	bytes, err := manifestsConfigMapYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/configmap.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func manifestsServiceYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/service.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}

func manifestsServiceYaml() (*asset, error) {
	bytes, err := manifestsServiceYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/service.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func manifestsClusterRoleYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/clusterrole.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func manifestsClusterRoleYaml() (*asset, error) {
	bytes, err := manifestsClusterRoleYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/clusterrole.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func manifestsClusterRoleBindingYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/clusterrolebinding.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func manifestsClusterRoleBindingYaml() (*asset, error) {
	bytes, err := manifestsClusterRoleBindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/clusterrolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func manifestsPersistentVolumeClaimYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../manifests/persistentvolumeclaim.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func manifestsPersistentVolumeClaimYaml() (*asset, error) {
	bytes, err := manifestsPersistentVolumeClaimYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../manifests/persistentvolumeclaim.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

//manifestsPersistentVolumeClaimYaml

func alamedarecommendationsCRDYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/CustomResourceDefinition/alamedarecommendationsCRD.yaml"), nil
}
func alamedarecommendationsCRDYaml() (*asset, error) {
	bytes, err := alamedarecommendationsCRDYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/CustomResourceDefinition/alamedarecommendationsCRD.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedascalersCRDYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/CustomResourceDefinition/alamedascalersCRD.yaml"), nil
}
func alamedascalersCRDYaml() (*asset, error) {
	bytes, err := alamedascalersCRDYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/CustomResourceDefinition/alamedascalersCRD.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func alamedadatahubCRBYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRoleBinding/alameda-datahubCRB.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func alamedadatahubCRBYaml() (*asset, error) {
	bytes, err := alamedadatahubCRBYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRoleBinding/alameda-datahubCRB.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaoperatorCRBYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRoleBinding/alameda-operatorCRB.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func alamedaoperatorCRBYaml() (*asset, error) {
	bytes, err := alamedaoperatorCRBYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRoleBinding/alameda-operatorCRB.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaevictionerCRBYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func alamedaevictionerCRBYaml() (*asset, error) {
	bytes, err := alamedaevictionerCRBYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRoleBinding/alameda-evictionerCRB.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func admissioncontrollerCRBYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRoleBinding/admission-controllerCRB.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func admissioncontrollerCRBYaml() (*asset, error) {
	bytes, err := admissioncontrollerCRBYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRoleBinding/admission-controllerCRB.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func alamedadatahubCRYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRole/alameda-datahubCR.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func alamedadatahubCRYaml() (*asset, error) {
	bytes, err := alamedadatahubCRYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRole/alameda-datahubCR.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaoperatorCRYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRole/alameda-operatorCR.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func alamedaoperatorCRYaml() (*asset, error) {
	bytes, err := alamedaoperatorCRYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRole/alameda-operatorCR.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaevictionerCRYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRole/alameda-evictionerCR.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func alamedaevictionerCRYaml() (*asset, error) {
	bytes, err := alamedaevictionerCRYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRole/alameda-evictionerCR.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func admissioncontrollerCRYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRole/admission-controllerCR.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func admissioncontrollerCRYaml() (*asset, error) {
	bytes, err := admissioncontrollerCRYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRole/admission-controllerCR.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func aggregatealamedaadmineditalamedaCRYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}
func aggregatealamedaadmineditalamedaCRYaml() (*asset, error) {
	bytes, err := aggregatealamedaadmineditalamedaCRYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ClusterRole/aggregate-alameda-admin-edit-alamedaCR.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func alamedadatahubSAYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ServiceAccount/alameda-datahubSA.yaml"), nil
	//return _manifestsServiceaccountYaml, nil
}

func alamedadatahubSAYaml() (*asset, error) {
	bytes, err := alamedadatahubSAYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ServiceAccount/alameda-datahubSA.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaoperatorSAYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ServiceAccount/alameda-operatorSA.yaml"), nil
	//return _manifestsServiceaccountYaml, nil
}

func alamedaoperatorSAYaml() (*asset, error) {
	bytes, err := alamedaoperatorSAYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ServiceAccount/alameda-operatorSA.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaevictionerSAYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ServiceAccount/alameda-evictionerSA.yaml"), nil
	//return _manifestsServiceaccountYaml, nil
}

func alamedaevictionerSAYaml() (*asset, error) {
	bytes, err := alamedaevictionerSAYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ServiceAccount/alameda-evictionerSA.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func admissioncontrollerSAYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ServiceAccount/admission-controllerSA.yaml"), nil
	//return _manifestsServiceaccountYaml, nil
}

func admissioncontrollerSAYaml() (*asset, error) {
	bytes, err := admissioncontrollerSAYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ServiceAccount/admission-controllerSA.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func grafanadatasourcesCMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/ConfigMap/grafana-datasources.yaml"), nil
	//return _manifestsSecretTokenYaml, nil
}

func grafanadatasourcesCMYaml() (*asset, error) {
	bytes, err := grafanadatasourcesCMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/ConfigMap/grafana-datasources.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func myalamedainfluxdbPVCYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml"), nil

}

func myalamedainfluxdbPVCYaml() (*asset, error) {
	bytes, err := myalamedainfluxdbPVCYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/PersistentVolumeClaim/my-alamedainfluxdbPVC.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func myalamedagrafanaPVCYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml"), nil

}

func myalamedagrafanaPVCYaml() (*asset, error) {
	bytes, err := myalamedagrafanaPVCYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/PersistentVolumeClaim/my-alamedagrafanaPVC.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func alamedadatahubSVYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Service/alameda-datahubSV.yaml"), nil

}

func alamedadatahubSVYaml() (*asset, error) {
	bytes, err := alamedadatahubSVYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Service/alameda-datahubSV.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func admissioncontrollerSVYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Service/admission-controllerSV.yaml"), nil

}

func admissioncontrollerSVYaml() (*asset, error) {
	bytes, err := admissioncontrollerSVYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Service/admission-controllerSV.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedainfluxdbSVYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Service/alameda-influxdbSV.yaml"), nil

}

func alamedainfluxdbSVYaml() (*asset, error) {
	bytes, err := alamedainfluxdbSVYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Service/alameda-influxdbSV.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedagrafanaSVYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Service/alameda-grafanaSV.yaml"), nil

}

func alamedagrafanaSVYaml() (*asset, error) {
	bytes, err := alamedagrafanaSVYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Service/alameda-grafanaSV.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

func alamedadatahubDMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Deployment/alameda-datahubDM.yaml"), nil

}

func alamedadatahubDMYaml() (*asset, error) {
	bytes, err := alamedadatahubDMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Deployment/alameda-datahubDM.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaoperatorDMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Deployment/alameda-operatorDM.yaml"), nil

}

func alamedaoperatorDMYaml() (*asset, error) {
	bytes, err := alamedaoperatorDMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Deployment/alameda-operatorDM.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedaevictionerDMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Deployment/alameda-evictionerDM.yaml"), nil

}

func alamedaevictionerDMYaml() (*asset, error) {
	bytes, err := alamedaevictionerDMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Deployment/alameda-evictionerDM.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func admissioncontrollerDMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Deployment/admission-controllerDM.yaml"), nil

}

func admissioncontrollerDMYaml() (*asset, error) {
	bytes, err := admissioncontrollerDMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Deployment/admission-controllerDM.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedainfluxdbDMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Deployment/alameda-influxdbDM.yaml"), nil

}

func alamedainfluxdbDMYaml() (*asset, error) {
	bytes, err := alamedainfluxdbDMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Deployment/alameda-influxdbDM.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func alamedagrafanaDMYamlBytes() ([]byte, error) {
	return ReadYamlFile("../../assets/Deployment/alameda-grafanaDM.yaml"), nil

}

func alamedagrafanaDMYaml() (*asset, error) {
	bytes, err := alamedagrafanaDMYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../assets/Deployment/alameda-grafanaDM.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

