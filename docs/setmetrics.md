# Metrics

**setp1**

install alameda-ai component

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
    team: frontend
spec:
  endpoints:
  - port: ai-metrics
  namespaceSelector:
    any: true
  selector:
    matchLabels:
      app: alameda-ai

```

**setp2**

Add your prometheus-k8s rbac

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