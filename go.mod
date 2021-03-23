module github.com/open-cluster-management/cert-policy-controller

go 1.14

require (
	github.com/onsi/gomega v1.8.1
	github.com/operator-framework/operator-sdk v0.17.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	golang.org/x/text => golang.org/x/text v0.3.3 // CVE-2020-14040
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)
