# QuickStart

The **FederatorAI-Operator** is an operator that manage [Alameda](https://github.com/containers-ai/alameda) in ways of:
- Deployment
- Upgrade
- Application Lifecycle and storage

And this document helps you to get started. In the following sections, we first show how to install **FederatorAI-Operator** and then how to use it.

## Deployment

Like any Kubernetes application, the deployment of a Kubernetes application can directly apply K8s manifests or leverage 3rd-party tools/frameworks. Here we provides but not limit to two ways:
- by K8s manifests  
- by operator-lifecycle-management framework

During the deployment, **FederatorAI-Operator** will install a CRD called _AlamedaService_ as a channel for users to interact with it. **FederatorAI-Operator** will reconcile to an _AlamedaService_ CR in a cluster wide scope.

#### Deployment by K8s Manifests

1. Checkout **FederatorAI-Operator** repository from github  
```
$ git clone https://github.com/containers-ai/federatorai-operator.git
```

2. Apply the K8s manifests
```
$ kubectl apply -f federatorai-operator/deploy/upstream/
```
This will pull **FederatorAI-Operator** image from `quay.io/prophetstor` and deploy it into `federatorai` namespace. You should see `federatorai-operator` pod is running after few seconds.

#### Deployment by Operator-Lifecycle-Management Framework

[Operator-Lifecycle-Management(OLM)](https://github.com/operator-framework/operator-lifecycle-manager) extends Kubernetes to provide a declarative way to install, manage, and upgrade operators and their dependencies in a cluster. To deploy **FederatorAI-Operator** by OLM, please follow the instructions at [OperatorHub.io](https://operatorhub.io/operator/alpha/federatorai.v0.0.1). Here copies the instructions as a quick reference.

1. Install OLM first
```
$ kubectl create -f https://raw.githubusercontent.com/operator-framework/operator-lifecycle-manager/master/deploy/upstream/quickstart/olm.yaml
```

2. Install **FederatorAI-Operator**
```
$ kubectl create -f https://operatorhub.io/install/alpha/federatorai.v0.0.1.yaml
```
This will pull image from `quay.io/prophetstor` and install **FederatorAI-Operator** version 0.0.1 to `operators` namespace. You should see `federatorai-operator` pod is running after few seconds.

## Using FederatorAI-Operator

To use **FederatorAI-Operator**, users need to create/apply an _AlamedaService_ CR in a namespace. Here is an example of _AlamedaService_ CR.
```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  enableexecution: true
  enablegui: true
  version: latest               ## for Alameda components. (exclude grafana and influxdb)
  prometheusservice: https://prometheus-k8s.openshift-monitoring:9091
  storages:                     ## see following details for where it is used in each component
    - usage: log                ## storage setting for log
      type: ephemeral           ## ephemeral means emptyDir{}
    - usage: data               ## storage setting for data
      type: pvc                 ## pvc means PersistentVolumeClaim
      size: 10Gi                ## mandatory when type=pvc
      class: "normal"           ## mandatory when type=pvc
```
By creating this CR, **FederatorAI Operator** starts to:
- deploy Alameda core components, components for recommendation execution and components for GUI to `alameda` namespace.
- The pulled Alameda component image tag is _latest_ except InfluxDB and Grafana components.
- Alameda datahub will retrieve metrics from Prometheus at _https://prometheus-k8s.openshift-monitoring:9091_
- log path will be mounted with _emptyDir{}_ for each component
- PVC are claimed and mounted in data path for each component

For more details, please refer to [AlamedaService CRD document](./crd_alamedaservice.md).


In addition, users can patch a created _AlamedaService_ CR and **FederatorAI-Operator** will react to it. For example, by changing the _enableexecution_ field from _true_ to _false_, Alameda recommendation execution components will be uninstalled. (Alameda is still giving prediction and recommendations. GUI can still visualize the result. Just the execution part is off)

