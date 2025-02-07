package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	v1 "open-cluster-management.io/cert-policy-controller/api/v1"
)

var (
	certPolicyStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificatepolicy_governance_info",
			Help: "The compliance status of the named CertificatePolicy. " +
				"0 == Compliant. 1 == NonCompliant. -1 == Unknown/Pending",
		},
		[]string{
			"policy",           // The name of the configuration policy
			"policy_namespace", // The namespace where the configuration policy is defined
		},
	)
	policyTotalEvalSecondsCounter = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "cert_policy_reconcile_seconds",
		Help: "The time in seconds it takes to evaluate all certificate policies.",
	})
	policyEvalSecondsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cert_policy_evaluation_seconds_total",
			Help: "The total seconds taken while evaluating the certificate policy. Use this alongside " +
				"cert_policy_evaluation_seconds_total.",
		},
		[]string{"name"},
	)
	policyEvalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cert_policy_evaluation_total",
			Help: "The total number of evaluations of the certificate policy. Use this alongside " +
				"cert_policy_evaluation_total.",
		},
		[]string{"name"},
	)
)

func init() {
	// Register custom metrics with the global Prometheus registry
	metrics.Registry.MustRegister(
		certPolicyStatusGauge,
		policyTotalEvalSecondsCounter,
		policyEvalSecondsCounter,
		policyEvalCounter,
	)
}

func getStatusValue(complianceState v1.ComplianceState) float64 {
	if complianceState == v1.Compliant {
		return 0
	} else if complianceState == v1.NonCompliant {
		return 1
	}

	return -1
}

func removeCertPolicyMetrics(request ctrl.Request) {
	// If a metric has an error while deleting, that means the policy was never evaluated so it can be ignored.
	_ = certPolicyStatusGauge.DeletePartialMatch(prometheus.Labels{
		"policy":           request.Name,
		"policy_namespace": request.Namespace,
	})
	_ = policyEvalSecondsCounter.DeleteLabelValues(request.Name)
	_ = policyEvalCounter.DeleteLabelValues(request.Name)
}
