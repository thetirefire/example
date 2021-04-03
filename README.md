# BadIdea Example project

```sh
cat << EOF > test.kubeconfig
apiVersion: v1
kind: Config
users:
- name: badidea
  user:
    username: bad
    password: idea
clusters:
- name: badidea
  cluster:
    server: https://localhost:6443
    insecure-skip-tls-verify: true
contexts:
- name: badidea
  context:
    cluster: badidea
    user: badidea
current-context: badidea
EOF
kustomize build config/crd | kubectl --kubeconfig=test.kubeconfig apply -f -
go run main.go --kubeconfig=test.kubeconfig
```