## AlamedaService Custom Resource Definition

**FederatorAI Operator** provides `alamedaservice` [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) as a channel for users to manage Alameda components. Here is an example of `alamedaservice` CR:

```
apiVersion: federatorai.containers.ai/v1alpha1
kind: AlamedaService
metadata:
  name: my-alamedaservice
  namespace: alameda
spec:
  enableexecution: true
  enablegui: true
  persistentvolumeclaim: ''
  prometheusservice: 'https://prometheus-k8s.openshift-monitoring:9091'
  version: latest
```

In this example, it creates an `AlamedaService` CR with name `my-alamedaservice` in namespace `alameda`. By creating this CR, **FederatorAI Operator** starts to:
- deploy Alameda core components, components of recommendation execution and components of GUI
- the container image version of deployed components are set to tag *latest*
- Alameda datahub will retrieve metrics from Prometheus at *https://prometheus-k8s.openshift-monitoring:9091*
- the **persistentcolumeclaim** also provides settings for mount a persistent volume for storage of prediction metrics, GUI dashboards and component logs.

Users can also modify an *AlamedaService* CR to alter the configurations. **FederatorAI Operator** will seemless update Alameda settings accordingly.

