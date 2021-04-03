module github.com/thetirefire/example

go 1.16

require (
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/thetirefire/badidea v0.0.0-20210403205638-26f5cc7e55b6
	k8s.io/apimachinery v0.21.0-beta.1
	k8s.io/client-go v0.21.0-beta.1
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.0-alpha.1
)

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
