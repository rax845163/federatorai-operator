
# <img src="./logo.png" width=60/> Federator.ai Operator<img src="./rhcert.png" align="right" width=150/>

**Federator.ai Operator** is an Operator that manages **Federator.ai** components for an OpenShift cluster. Once installed, the Federator.ai Operator provides the following features:
- **Create/Clean up**: Launch **Federator.ai** components using the Operator.
- **Easy Configuration**: Easily configure data source of Prometheus and enable/disable add-on components, such as GUI, and predictive autoscaling.
- **Pod Scaling Recommendation/Autoscaling**: Use provided CRD to setup target pods and desired policies for scaling recommendation and autoscaling.

> **Note:** **Federator.ai** requires a Prometheus datasource to get historical metrics of pods and nodes. When launching **Federator.ai** components, Prometheus connection settings need to be provided.

## Federator.ai

**Federator.ai** is the brain of resource orchestration for kubernetes. We use machine learning technology to provide intelligence that foresees future resource usage of your Kubernetes cluster across multiple layers. Federator.ai recommends the right sizes of containers and the right number of replications. It also elastically manages pod scaling and scheduling of your containerized applications. The overall benefit is cost reduction up to 50% and higher service quality, such as fewer OOM issues. For more information, visit [github.com/containers-ai/alameda](https://github.com/containers-ai/alameda)

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
