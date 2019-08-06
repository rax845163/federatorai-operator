## AlamedaService Custom Resource Definition

**Federator.ai Operator** provides _AlamedaService_ [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) as a channel for users to manage Alameda components including:
- Deployment of Alameda components such as _alameda-operator_, _alameda-datahub_, _alameda-ai_, _alameda-evictioner_, _alameda-admission-controller_, _alameda-recommender_, _InfluxDB_ and _Grafana_. Please visit [Alamede architecture](https://github.com/containers-ai/alameda/blob/master/design/architecture.md) for more details.
- Seamless updation of Alameda between versions.
- Application lifecycle and storage management.

An _AlamedaService_ CR is structured as:
- a section of global setting
  This section provides a configurable global setting for all Alameda components. This section intends to provide a fast and easy setup to quickstart Alameda.
- a section of detailed setting for each component (optional)
  The settings in this section are optional and it is used to fine tune the values inherited from the global section for each Alameda component.

When an _AlamedaService_ CR is created, Federator.ai Operator will reconcile it and spawn operands. For the detail schema of _AlamedaService_, please refer to the last section of this document. Here we shows two example to quickly give users a feel what the configuration that an _AlamedaService_ tries to provide.

### An Example
Here is an _AlamedaService_ CR example:

```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  selfDriving: true             ## to enable resource self-orchestration of the deployed Alameda components
                                ## it is recommended NOT to use ephemeral data storage for Alameda influxdb component when self-Driving is enabled
  enableExecution: true
  enableGui: true
  version: latest               ## for Alameda components. (exclude grafana and influxdb)
  prometheusService: https://prometheus-k8s.openshift-monitoring:9091
  storages:                     ## see following details for where it is used in each component
    - usage: log                ## storage setting for log
      type: ephemeral           ## ephemeral means emptyDir{}
    - usage: data               ## storage setting for data
      type: pvc                 ## pvc means PersistentVolumeClaim
      size: 10Gi                ## mandatory when type=pvc
      class: "normal"           ## mandatory when type=pvc
```

In this example, it creates an _AlamedaService_ CR with name _my-alamedaservice_ in namespace `alameda`. By creating this CR, **Federator.ai Operator** starts to:
- deploy Alameda core components, components for recommendation execution and components for GUI
- create an [_AlamedaScaler_](https://github.com/containers-ai/alameda/blob/master/design/crd_alamedascaler.md) to self-orchestrate Alameda's resource usage
- pull _latest_ Alameda component image except InfluxDB and Grafana components. To overwrite the pulled image tag of InfluxDB and Grafana, users can specify them in _section schema for each component_.
- set Alameda datahub to retrieve metrics from Prometheus at _https://prometheus-k8s.openshift-monitoring:9091_
- mount _emptyDir{}_ to log path for each component
- claim volumn by PVC and mount it to data path for each component

### A More Complicated Example
Here is another _AlamedaService_ CR example to show how to overwrite the global setting for some components:

```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  selfDriving: true             ## to enable resource self-orchestration of the deployed Alameda components
                                ## it is recommended NOT to use ephemeral data storage for Alameda influxdb component when self-Driving is enabled
  enableExecution: true
  enableGui: true
  version: v0.3.38              ## for Alameda components. (exclude grafana and influxdb)
  prometheusService: https://prometheus-k8s.openshift-monitoring:9091
  storages:
    - usage: log                ## storage for log of each component
      type: ephemeral
    - usage: data               ## storage for data of each component
      type: pvc
      size: 10Gi
      class: "normal"

# following are more detail configurations for each component and overwrite the global config
# for complete list of Alameda components, please visit https://github.com/containers-ai/federatorai-operator/blob/master/docs/crd_alamedaservice.md
  alameda-ai:
    image: quay.io/prophetstor/alameda-ai
    version: latest
    imagePullPolicy: Always
    storages:
      usage: log      ## for path /var/log/alameda
      type: pvc
      size: 10Gi
      class: "normal"

  alameda-grafana:
    image: grafana/grafana
    version: 5.4.3
    storages:
      usage: data     ## for path /var/lib/grafana
      type: pvc
      size: 1Gi
      class: "normal"

  alameda-influxdb:
    image: influxdb
    version: 1.7-alpine
    storages:
      usage: data     ## for path /var/lib/influxdb
      type: pvc
      size: 20Gi
      class: "fast"
```

## Schema of AlamedaService

- Field: metadata
  - type: ObjectMeta
  - description: This follows the ObjectMeta definition in [Kubernetes API Reference](https://kubernetes.io/docs/reference/#api-reference).
- Field: spec
  - type: [AlamedaServiceSpec](#alamedaservicespec)
  - description: Spec of AlamedaService.

### AlamedaServiceSpec

- Field: selfDriving
  - type: boolean
  - description: If this field is set to _true_, Federator.ai Operator will create an [_AlamedaScaler_](https://github.com/containers-ai/alameda/blob/master/design/crd_alamedascaler.md) CR to self-orchestrate the resource usage of deployed Alameda components. Default is _false_.
> **Note:** It is highly recommended to use persistent storage for data in Alameda influxdb component when self-Driving is enabled
- Field: platform
  - type: string
  - description: (Optional) Specify this key with "openshift3.9" value if you are deploying Alameda in OCP/OKD 3.9 environment or the execution of container cpu and memory resource limit/request changes may not work. For other platforms, this key is optional.
- Field: enableExecution
  - type: boolean
  - description: Federator.ai Operator will deploy components to execute _AlamedaRecommendation_ CRs if this field is set to _true_. Default is _false_.
- Field: enableGui
  - type: boolean
  - description: Federator.ai Operator will deploy GUI to visualize Alameda predictions/recommendations and cluster/node status if this field is set to _true_. Default is _true_.
- Field: enableFedemeter
  - type: boolean
  - description: Federator.ai Operator will deploy Fedemeter and you must add your Fedemeter serviceAccount to privileged SecurityContextConstraints
- Field: version
  - type: string
  - description: It sets the version tag when pulling Alameda component images.
- Field: prometheusService
  - type: string
  - description: This field tells datahub and Grafana where the Prometheus URL is to retrieve pods/nodes peformance metrics data.
- Field: storages
  - type: [StorageSpec](#storagespec) array
  - description: This field is optional and it lists storage settings which applied to each operand.
- Field: serviceExposures
  - type: [ServiceExposureSpec](#serviceexposurespec) array
  - description: This field is optional and it lists service exposure settings which applied to an Alameda component.
- Field: alamedaOperator
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-operator component. This field is optional.
- Field: alamedaDatahub
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-datahub component. This field is optional.
- Field: alamedaAi
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-ai component. This field is optional.
- Field: alamedaEvictioner
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-evictioner component. This field is optional.
- Field: alamedaAdmissionController
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-admission-controller component. This field is optional.
- Field: alamedaInfluxdb
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for InfluxDB component. This field is optional.
- Field: alamedaGrafana
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-grafana component. This field is optional.
- Field: alamedaRecommender
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-recommender component. This field is optional.
- Field: alamedaExecutor
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-executor component. This field is optional.
- Field: fedemeter
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for fedemeter component. This field is optional.

### StorageSpec

- Field: usage
  - type: string
  - description: This field is mandatory and the supported values are _log_ and _data_.
- Field: type
  - type: string
  - description: The supported values of this field are _ephemeral_ and _pvc_. _ephemeral_ means this storage will be mounted with [_emptyDir{}_](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) and _pvc_ means this storage will be allocated with [_PersistentVolumeClaim_](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims). The default value is _ephemeral_.
- Field: size
  - type: object
  - description: This field applies to _pvc_ and _ephemeral_ type. It claims a persistent volume from K8s with the size and is mandatory if type is _pvc_. For how to setup the value, visit [capacity](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes) setting of a Kubernetes persistent volume.
- Field: class
  - type: string
  - description: This field is mandatory if type is _pvc_. It claims a persistent volume from K8s with the matching _storageClassName_.
- Field: accessMode
  - type: string
  - description: This field is for _pvc_ type. Is sets how the claimed volume is mounted. Default is _ReadWriteOnce_.

### AlamedaComponentSpec

- Field: image
  - type: string
  - description: It sets the image name to be pulled.
- Field: imagePullPolicy
  - type: string
  - description: It follows the K8s [image spec](https://kubernetes.io/docs/concepts/containers/images/) of container to pull images. Default is _IfNotPresent_.
- Field: version
  - type: string
  - description: It sets the version tag of image to be pulled.
- Field: storages
  - type: [StorageSpec](#storagespec) array
  - description: This field is optional and it lists storage settings which applied to an Alameda component.

### ServiceExposureSpec

- Field: name
  - type: string
  - description: It sets the service name to be exposed. The list of available names are equal to the services' name under [folder](../assets/Service).
- Field: type
  - type: string
  - description: It sets the type of service exposure. Currently supported type is **NodePort**.
- Field: nodePort
  - type: [NodePortSpec](#nodeportspec)
  - description: This field will be applied when type is NodePort.

### NodePortSpec

- Field: ports
  - type: [PortSpec](#portspec) array
  - description: This field lists the ports to be exposed by NodePort type.

### PortSpec

- Field: port
  - type: integer
  - description: It sets which service port to be proxied. 
- Field: nodePort
  - type: integer
  - description: It sets which port on the node to proxy to the service port.







