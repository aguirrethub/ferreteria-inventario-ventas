package http_handlers

import (
	"encoding/json"
	"net/http"

	"ferreteria-inventario-ventas/internal/service"
)

// Handlers agrupa los servicios.
type Handlers struct {
	ClientsSvc  *service.ClientService
	ProductsSvc *service.ProductService
	SalesSvc    *service.SaleService
}

// Función auxiliar para responder JSON.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Health verifica que el servidor está funcionando.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{
		"status": "ok",
	})
}
