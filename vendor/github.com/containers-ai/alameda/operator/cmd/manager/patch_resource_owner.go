package main

import (
	"context"
	"time"

	utils "github.com/containers-ai/alameda/pkg/utils"
	k8sUtils "github.com/containers-ai/alameda/pkg/utils/kubernetes"
	"github.com/containers-ai/alameda/pkg/utils/kubernetes/metadata"
	openshiftappsv1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func addOwnerReferenceToResourcesCreateFrom3rdPkg(sigsK8SClient client.Client) {

	dep, dc, err := getDeploymentOrDeploymentConfigRunningManager()
	if err != nil {
		scope.Errorf("add owner reference to resources create from 3rd pkg failed: %s", err.Error())
		return
	} else if dep == nil && dc == nil {
		scope.Error("add owner reference to resources create from 3rd pkg failed, cannot get deployment or deploymentConfig owning alameda-operator")
		return
	}

	var ownerType metav1.TypeMeta
	var ownerMeta metav1.ObjectMeta
	if dep != nil {
		ownerType.APIVersion = appsv1.SchemeGroupVersion.String()
		ownerType.Kind = "Deployment"
		ownerMeta = dep.ObjectMeta
	} else if dc != nil {
		ownerType.APIVersion = openshiftappsv1.SchemeGroupVersion.String()
		ownerType.Kind = "DeploymentConfig"
		ownerMeta = dc.ObjectMeta
	}
	ownerRef := k8sUtils.NewOwnerReference(ownerType, ownerMeta, false)

	retryPeriod := 60 * time.Second
	retryTicker := time.NewTicker(retryPeriod)
	for range retryTicker.C {

		retry := false

		serviceKeys := []client.ObjectKey{client.ObjectKey{Namespace: operatorConf.K8SWebhookServer.Service.Namespace, Name: operatorConf.K8SWebhookServer.Service.Name}}
		for _, servicesKey := range serviceKeys {
			service := corev1.Service{}
			if err := sigsK8SClient.Get(context.TODO(), servicesKey, &service); err != nil {
				scope.Warnf("add ownerReferences to service: %s/%s failed, retry after %f seconds, %s.", servicesKey.Namespace, servicesKey.Name, retryPeriod.Seconds(), err.Error())
				retry = true
				break
			}
			k8sUtils.AddOwnerRefToObject(&service, ownerRef)
			if err := sigsK8SClient.Update(context.TODO(), &service); err != nil {
				scope.Warnf("add ownerReferences to service: %s/%s failed, retry after %f seconds, %s.", servicesKey.Namespace, servicesKey.Name, retryPeriod.Seconds(), err.Error())
				retry = true
				break
			}
		}

		validatingWebhookConfigKeys := []client.ObjectKey{client.ObjectKey{Name: operatorConf.K8SWebhookServer.ValidatingWebhookConfigName}}
		for _, webhookConfigKey := range validatingWebhookConfigKeys {
			webhookConfig := admissionregistrationv1beta1.ValidatingWebhookConfiguration{}
			if err := sigsK8SClient.Get(context.TODO(), webhookConfigKey, &webhookConfig); err != nil {
				scope.Warnf("add ownerReferences to validatingWebhookConfiguration: %s failed, retry after %f seconds, %s", webhookConfigKey.Name, retryPeriod.Seconds(), err.Error())
				retry = true
				break
			}
			k8sUtils.AddOwnerRefToObject(&webhookConfig, ownerRef)
			if err := sigsK8SClient.Update(context.TODO(), &webhookConfig); err != nil {
				scope.Warnf("add ownerReferences to validatingWebhookConfiguration: %s failed, retry after %f seconds, %s", webhookConfig.Name, retryPeriod.Seconds(), err.Error())
				retry = true
				break
			}
		}

		// mutatingWebhookConfigKeys := []client.ObjectKey{client.ObjectKey{Name: operatorConf.K8SWebhookServer.MutatingWebhookConfigName}}
		// for _, webhookConfigKey := range mutatingWebhookConfigKeys {
		// 	webhookConfig := admissionregistrationv1beta1.MutatingWebhookConfiguration{}
		// 	if err := sigsK8SClient.Get(context.TODO(), webhookConfigKey, &webhookConfig); err != nil {
		// 		scope.Errorf("add ownerReferences to mutatingWebhookConfiguration: %s failed, retry after %f seconds, %s", webhookConfig.Name,retryPeriod.Seconds(), err.Error())
		// 		retry = true
		// 		break
		// 	}
		// 	k8sUtils.AddOwnerRefToObject(&webhookConfig, ownerRef)
		// 	if err := sigsK8SClient.Update(context.TODO(), &webhookConfig); err != nil {
		// 		scope.Errorf("add ownerReferences to mutatingWebhookConfiguration: %s failed, retry after %f seconds, %s", webhookConfig.Name,retryPeriod.Seconds(), err.Error())
		// 		retry = true
		// 		break
		// 	}
		// }

		if !retry {
			scope.Info("add owner reference to resources create from 3rd pkg success")
			return
		}
	}
}

func getDeploymentOrDeploymentConfigRunningManager() (*appsv1.Deployment, *openshiftappsv1.DeploymentConfig, error) {

	podNamespace := utils.GetRunningNamespace()
	podName := utils.GetRunnningPodName()
	pod, err := k8sUtils.GetPodByNamespaceNameWithConfig(podNamespace, podName, *k8sConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get deployment or deploymentConfig running alameda-operator failed")
	}

	ort, err := metadata.NewOwnerReferenceTracerWithConfig(*k8sConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get deployment or deploymentConfig running alameda-operator failed")
	}
	dep, dc, err := ort.GetDeploymentOrDeploymentConfigOwningPod(pod)
	if err != nil {
		err = errors.Wrap(err, "get deployment or deploymentConfig running alameda-operator failed")
	}

	return dep, dc, err
}
