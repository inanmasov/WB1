package router

import (
	"net/http"

	handler "example.com/service/service/internal/transport/handler"
)

func RegisterRoutes() {
	http.HandleFunc("/", handler.HealthCheckHandler)
	http.HandleFunc("/order", handler.OrdersHandler)
	http.HandleFunc("/view", handler.OrderView)
}
