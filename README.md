# FederatorAI Operator

**FederatorAI Operator** is an Operator that manages **Alameda** components for an openshift cluster. Once installed, the FederatoAI Operator provides the following features:
- **Create/Clean up**: Launch **Alameda** components using the Operator.
- **Easy Configuration**: Easily config source of Prometheus and enable/disable addon components such as GUI, and execition.
- **Autoscaling Pod**: Use provided CRD to setup target pods for autoscaling.

> **Note:** **Alameda** requires a Prometheus datasource to get historical metrics of pods and nodes. When launching **Alameda** components, Prometheus connection settings need to be provided.

## Alameda

**Alameda** is the brain of resource orchestration for kubernetes. It foresees future resource usage of your Kubernetes cluster from the cloud layer down to the pod level. We use machine learning technology to provide intelligence that enables dynamic scaling and scheduling of your containers - effectively making us the “brain” of Kubernetes resource orchestration. By providing full foresight of resource availability, demand, health, impact and SLA, we enable cloud strategies that involve changing provisioned resources in real time. For more information, visit [github.com/containers-ai/alameda](https://github.com/containers-ai/alameda)

## Documentations
Please visit [docs](./docs/)

## Resources

* [How to Write Go Code](https://golang.org/doc/code.html)
* [Effective Go](https://golang.org/doc/effective_go.html)
* [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* [GitHub Flow](https://guides.github.com/introduction/flow/)
* [Go modules](https://github.com/golang/go/wiki/Modules)
* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [CRI tools](https://github.com/kubernetes-sigs/cri-tools)
