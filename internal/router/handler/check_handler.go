package handler

import (
	"fmt"
	"net/http"

	"github.com/nickyrolly/ws-chat-demo/internal/domain"
	"github.com/nickyrolly/ws-chat-demo/internal/usecase"
)

func CheckServices(w http.ResponseWriter, r *http.Request) {
	healthCheckResult := make(map[string]string)
	for name, host := range domain.ServiceMap {
		usecase.ServiceHealthCheck(&healthCheckResult, name, host)
	}
	usecase.PublishHealthCheck(&healthCheckResult)

	for serviceName, status := range healthCheckResult {
		fmt.Fprintf(w, "%s: %s\n", serviceName, status)
	}
}
