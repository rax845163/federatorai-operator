## Deploy Admission-Controller

1. Deploy Secret and MutatingWebhookConfiguration that admission-controller needs.
```
sh $GOPATH/src/github.com/containers-ai/alameda/admission-controller/generate-key.sh
```
2. Deploy admission-controller.
```
kubectl apply -f $GOPATH/src/github.com/containers-ai/alameda/admission-controller/deployment.yaml
```