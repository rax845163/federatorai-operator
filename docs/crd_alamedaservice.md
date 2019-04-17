## AlamedaService Custom Resource Definition

**FederatorAI Operator** provides `alamedaservice` [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) as a channel for users to manage Alameda components. The structure of an `alamedaservice` CR spec is composed of
- a section of global setting
  This section provides a configurable global setting for all Alameda components. This section intends to provide a fast and easy setup to quickstart Alameda.
- a section of setting for each component
  This sections provides a more comprehensive settings of alameda components and overwrite configure from global settings if they are conflict.

Here is a depiction of an `alamedaservice` CR:
```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  <global section schema>
  <section schema for each component>
```
When this `alamedaservice` CR is applied, FederatorAI-Operator will reconcile it and spawn operands accordingly.
The following explains the *global section* and *section for each component*

### Global Section Schema

Here lists the global section schema:
- enableexecution  
  This field is boolean type. FederatorAI Operator will deploy components to execute AlamedaRecommendation if enabled.
- enablegui  
  This field is boolean type. FederatorAI Operator will deploy GUI to visualize AlamedaRecommnedation if enabled.
- imageHub  
  This field is string type. It sets the image hub where the Alameda components are pull from.
- imageBase  
  This field is string type. It sets the base images what the Alameda components are on top of. Currently Alameda may built from rhel, ubi and alpine image bases.
- imagePullPolicy  
  This field is string type. It follows the K8s [image spec](https://kubernetes.io/docs/concepts/containers/images/) of container to pull images. Default is *IfNotPresent*
- version  
  This field is string type. It sets the version tag of Alameda component images.
> **Note:** The image location to be pulled is devided by field *imageHub*, *imageBase* and version with format "*imageHub*/<component name>-*imageBase*:*version*"
- prometheusservice  
  This field is string type. It tells datahub and Grafana where the Prometheus is to retrieve pods/nodes peformance metrics data.
- storages  
  This field is a list of storage schema. If provides volume settings which apply to each Alameda component.

And here lists the storage schema:
- usage  
  This field is string type and right now *log* and *data* is supported. *log* means this storage setting is for log and */var/log/<application>* will be mounted. *data* means */var/lib/<application>* will be mounted. The default value is "" which means it applied to all possible value in this field (which are *log* and *data*).
- type  
  This field is string type and right now *ephemeral* and *pvc* is supported. *ephemeral* means this storage will use *emptyDir{}* and *pvc* means this *PersistentVolumeClaim*. The default value is *ephemeral*.
- size  
  This field is integer type and is mandatory if type is *pvc*. It claims a persistent volume from K8s with the size.
- class  
  This field is string type and is mandatory if type is *pvc*. It claims a persistent volume from K8s of the class.
- accessMode  
  This field is string type. Is shows how the claimed volume is mounted. Default is *ReadWriteOnce*

### Section Schema for Each Component

Here lists the schema for each component:
- image  
  This field is a string type. It sets the image name to be pulled.
- imagePullPolicy  
  This field is string type. It follows the K8s [image spec](https://kubernetes.io/docs/concepts/containers/images/) of container to pull images.
- version  
  This field is a string type. It sets the version tag of image to be pulled.
> **Note:** The image location to be pulled is devided by field *image* and version with format "*image*:*version*"
- storages  
  This is also the storage schema described above.
  

### Example
Here is an example of `alamedaservice` CR:

```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  enableexecution: true
  enablegui: true
  imageHub: quay.io/prophetstor ## for alameda components. (exclude grafana and influxdb)
  imageBase: ubi                ## for alameda components. (exclude grafana and influxdb)
  version: latest               ## for alameda components. (exclude grafana and influxdb)
  prometheusservice: https://prometheus-k8s.openshift-monitoring:9091
  storages:                     ## see following details for where it is used in each component
    - usage: log                ## storage for log, or data
      type: ephemeral           ## default is ephemeral means emptyDir. pvc means persistent volume claim
    - usage: data               ## storage for log, or data
      type: pvc                 ## default is ephemeral means emptyDir. pvc means persistent volume claim
      size: 10Gi                ## mandatory if type=pvc;
      class: "normal"           ## mandatory if type=pvc;
```

In this example, it creates an `AlamedaService` CR with name `my-alamedaservice` in namespace `alameda`. By creating this CR, **FederatorAI Operator** starts to:
- deploy Alameda core components, components of recommendation execution and components of GUI
- The images of Alameda components are pull from *quay.io* with latest tag.
- Alameda datahub will retrieve metrics from Prometheus at *https://prometheus-k8s.openshift-monitoring:9091*
- log path will be mounted with *emptyDir*
- PVC are claimed and mounted for data path

### A More Complicated Example
Here is another example of `alamedaservice` CR to show how to overwrite the global setting for some components:

```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  enableexecution: true
  enablegui: true
  imageHub: quay.io/prophetstor ## for alameda components. (exclude grafana and influxdb)
  imageBase: ubi                ## for alameda components. (exclude grafana and influxdb)
  version: v0.3.7               ## for alameda components. (exclude grafana and influxdb)
  prometheusservice: https://prometheus-k8s.openshift-monitoring:9091
  storages:                     ## see following details for where it is used in each component
    - usage: log                ## storage for log, or data
      type: ephemeral           ## default is ephemeral means emptyDir. pvc means persistent volume claim
    - usage: data               ## storage for log, or data
      type: pvc                 ## default is ephemeral means emptyDir. pvc means persistent volume claim
      size: 10Gi                ## mandatory if type=pvc;
      class: "normal"           ## mandatory if type=pvc;

# following is more detail configures for each components and overwrite the global config
# components are: alameda-operator, alameda-datahub, alameda-ai, alameda-evictioner,
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

