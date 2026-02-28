package http_handlers

import (
	"encoding/json"
	"net/http"

	"ferreteria-inventario-ventas/internal/domain"
)

// Clients maneja:
// GET  -> listar clientes
// POST -> crear cliente

// Clients godoc
// @Summary Listar o crear clientes
// @Description GET lista clientes, POST crea cliente
// @Tags Clients
// @Accept json
// @Produce json
// @Param client body domain.Client false "Cliente (solo POST)"
// @Success 200 {array} domain.Client
// @Success 201 {object} domain.Client
// @Router /clients [get]
// @Router /clients [post]
func (h *Handlers) Clients(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		list, err := h.ClientsSvc.List()
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, list)

	case http.MethodPost:
		var input domain.Client

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "JSON inválido"})
			return
		}

		err = h.ClientsSvc.Create(&input)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 201, input)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
