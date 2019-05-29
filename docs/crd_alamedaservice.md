## AlamedaService Custom Resource Definition

**FederatorAI Operator** provides _AlamedaService_ [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) as a channel for users to manage Alameda components including:
- Deployment of Alameda components such as _alameda-operator_, _alameda-datahub_, _alameda-ai_, _alameda-evictioner_, _alameda-admission-controller_, _InfluxDB_ and _Grafana_. Please visit [Alamede architecture](https://github.com/containers-ai/alameda/blob/master/design/architecture.md) for more details.
- Seamless updation of Alameda between versions.
- Application lifecycle and storage management.

An _AlamedaService_ CR is structured as:
- a section of global setting  
  This section provides a configurable global setting for all Alameda components. This section intends to provide a fast and easy setup to quickstart Alameda.
- a section of detailed setting for each component (optional)  
  The settings in this section are optional and it is used to fine tune the values inherited from the global section for each Alameda component.

When an _AlamedaService_ CR is created, FederatorAI-Operator will reconcile it and spawn operands. For the detail schema of _AlamedaService_, please refer to the last section of this document. Here we shows two example to quickly give users a feel what the configuration that an _AlamedaService_ tries to provide.

### An Example
Here is an _AlamedaService_ CR example:

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

In this example, it creates an _AlamedaService_ CR with name _my-alamedaservice_ in namespace `alameda`. By creating this CR, **FederatorAI Operator** starts to:
- deploy Alameda core components, components for recommendation execution and components for GUI
- The pulled Alameda component image tag is _latest_. The only exceptions are InfluxDB and Grafana components. To overwrite the pulled image tag of InfluxDB and Grafana, users can specify them in _section schema for each component_.
- Alameda datahub will retrieve metrics from Prometheus at _https://prometheus-k8s.openshift-monitoring:9091_
- log path will be mounted with _emptyDir{}_ for each component
- PVC are claimed and mounted in data path for each component

### A More Complicated Example
Here is another _AlamedaService_ CR example to show how to overwrite the global setting for some components:

```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  enableexecution: true
  enablegui: true
  version: v0.3.7               ## for Alameda components. (exclude grafana and influxdb)
  prometheusservice: https://prometheus-k8s.openshift-monitoring:9091
  storages:
    - usage: log                ## storage for log of each component
      type: ephemeral
    - usage: data               ## storage for data of each component
      type: pvc
      size: 10Gi
      class: "normal"

# following are more detail configurations for each component and overwrite the global config
# Alameda components are: alameda-operator, alameda-datahub, alameda-ai, alameda-evictioner,
# alameda-admission-controller, alameda-grafana and alameda-influxdb
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

- Field: enableexecution
  - type: boolean
  - description: FederatorAI Operator will deploy components to execute AlamedaRecommendation CRs if this field is set to _true_. Default is _false_.
- Field: enablegui
  - type: boolean
  - description: FederatorAI Operator will deploy GUI to visualize Alameda predictions/recommendations and cluster/node status if this field is set to _true_. Default is _true_.
- Field: version
  - type: string
  - description: It sets the version tag when pulling Alameda component images.
- Field: prometheusservice
  - type: string
  - description: This field tells datahub and Grafana where the Prometheus URL is to retrieve pods/nodes peformance metrics data.
- Field: storages
  - type: [StorageSpec](#storagespec) array
  - description: This field is optional and it lists storage settings which applied to each operand.
- Field: alameda-operator
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-operator component. This field is optional.
- Field: alameda-datahub
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-datahub component. This field is optional.
- Field: alameda-ai
  - type: [AlamedaComponentSpec](#alamedacomponentspec) 
  - description: Spec for Alameda-ai component. This field is optional.
- Field: alameda-evictioner
  - type: [AlamedaComponentSpec](#alamedacomponentspec) 
  - description: Spec for Alameda-evictioner component. This field is optional.
- Field: alameda-admission-controller
  - type: [AlamedaComponentSpec](#alamedacomponentspec) 
  - description: Spec for Alameda-admission-controller component. This field is optional.
- Field: alameda-recommender
  - type: [AlamedaComponentSpec](#alamedacomponentspec) 
  - description: Spec for Alameda-recommender component. This field is optional.
- Field: alameda-influxdb
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for InfluxDB component. This field is optional.
- Field: alameda-grafana
  - type: [AlamedaComponentSpec](#alamedacomponentspec)
  - description: Spec for Alameda-grafana component. This field is optional.

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
- Field: accessmode
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












