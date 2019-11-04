# Metrics

**step1**

install alameda-ai and alameda-ai-dispatcher components

**setp2**

alameda-ai-servicemonitoring-cr.yaml
```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: alameda-ai-metrics
  namespace: openshift-monitoring
  labels:
    k8s-app: prometheus-operator
spec:
  endpoints:
  - port: ai-metrics
  namespaceSelector:
    any: true
  selector:
    matchLabels:
      component: alameda-ai
```

alameda-ai-dispatcher-servicemonitoring-cr.yaml
```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: alameda-ai-dispatcher-metrics
  namespace: openshift-monitoring
  labels:
    k8s-app: prometheus-operator
spec:
  endpoints:
  - port: ai-dispatcher-metrics
  namespaceSelector:
    any: true
  selector:
    matchLabels:
      component: alameda-ai-dispatcher
```

**step2**

Update your clusterrole prometheus-k8s rbac

```

- apiGroups:
  - ""
  attributeRestrictions: null
  resources:
  - endpoints
  - pods
  - services
  verbs:
  - list


```
