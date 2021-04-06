# BadIdea Example project

```sh
go run main.go

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

export KUBECONFIG=test.kubeconfig

kubectl api-resources

kubectl explain bar --recursive

kubectl create -f - <<EOF
---
apiVersion: foo.example.thetirefire/v1
kind: Bar
metadata:
  name: test
spec:
  color: blue
  shape: circle
EOF
```