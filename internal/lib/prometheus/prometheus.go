package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "go_exchange_auth",
			Name:      "login_attempts_total",
			Help:      "total number of login attempts",
		}, []string{"status"},
	)

	LoginErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "go_exchange_auth",
			Name:      "login_errors_total",
			Help:      "total number of failed login attempts",
		}, []string{"error_type"},
	)

	RegisterAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "go_exchange_auth",
			Name:      "register_attempts_total",
			Help:      "total number of registration attempts",
		}, []string{"status"},
	)

	RegisterErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "go_exchange_auth",
			Name:      "register_errors_total",
			Help:      "total number of failed register attempts",
		}, []string{"error_type"},
	)
)
