# BadIdea Example project

```sh
export KUBECTL="kubectl --server https://localhost:6443 --insecure-skip-tls-verify --username=bad --password=idea"
kustomize build config/crd | $KUBECTL apply -f -
make run
```