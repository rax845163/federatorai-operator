package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/containers-ai/alameda/operator/pkg/apis"
	appsapi "github.com/openshift/api/apps"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s_serializer_json "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func registerThirdPartyCRD() {
	apis.AddToSchemes = append(apis.AddToSchemes, appsapi.Install)
}

func applyCRDs(cfg *rest.Config) {
	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	crdFiles := []string{}
	if files, err := ioutil.ReadDir(crdLocation); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				crdFiles = append(crdFiles, crdLocation+string(os.PathSeparator)+file.Name())
			}
		}
	} else {
		scope.Error("Failed to read CRDs: " + err.Error())
	}

	for _, crdFile := range crdFiles {
		yamlBin, rfErr := ioutil.ReadFile(crdFile)
		if rfErr != nil {
			scope.Errorf(fmt.Sprintf("Read crd file %s failed.", crdFile))
			continue
		}

		s := k8s_serializer_json.NewYAMLSerializer(k8s_serializer_json.DefaultMetaFactory, scheme.Scheme,
			scheme.Scheme)

		var crdIns apiextensionsv1beta1.CustomResourceDefinition
		_, _, decErr := s.Decode(yamlBin, nil, &crdIns)
		if decErr != nil {
			scope.Errorf(fmt.Sprintf("Decode crd file %s failed: %s", crdFile, decErr.Error()))
			continue
		}

		_, createErr := apiextensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&crdIns)
		if createErr != nil {
			scope.Errorf(fmt.Sprintf("Failed to create CRD %s: %s", crdFile, createErr.Error()))
			_, updateErr := apiextensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Update(&crdIns)
			if updateErr != nil {
				scope.Errorf(fmt.Sprintf("Failed to update CRD %s: %s", crdFile, updateErr.Error()))
				continue
			}
		}
		err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
			crd, getErr := apiextensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdIns.Name, metav1.GetOptions{})
			if getErr != nil {
				scope.Warnf(fmt.Sprintf("Failed to wait for CRD %s creation: %s", crdFile, getErr.Error()))
				return false, nil
			}
			for _, cond := range crd.Status.Conditions {
				switch cond.Type {
				case apiextensionsv1beta1.Established:
					if cond.Status == apiextensionsv1beta1.ConditionTrue {
						scope.Infof(fmt.Sprintf("CRD %s created.", crdIns.Name))
						return true, nil
					}
				case apiextensionsv1beta1.NamesAccepted:
					if cond.Status == apiextensionsv1beta1.ConditionFalse {
						scope.Errorf(fmt.Sprintf("CRD name conflict: %v, %v", cond.Reason, err))
					}
				}
			}
			return false, nil
		})
		if err != nil {
			scope.Errorf(fmt.Sprintf("Polling crd for %s failed: %s", crdFile, err.Error()))
		}
	}
}
