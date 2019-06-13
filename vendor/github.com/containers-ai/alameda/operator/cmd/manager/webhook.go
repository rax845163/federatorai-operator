package main

import (
	"github.com/containers-ai/alameda/operator"
	k8swhsrv "github.com/containers-ai/alameda/operator/k8s-webhook-server"
	"github.com/containers-ai/alameda/operator/podinfo"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func launchWebhook(mgr *manager.Manager, config *operator.Config) {
	podInfo := podinfo.NewPodInfo(config.PodInfo)
	k8sWebhookSrv := k8swhsrv.NewK8SWebhookServer(mgr, config.K8SWebhookServer, podInfo.Labels)
	k8sWebhookSrv.Launch()
}
