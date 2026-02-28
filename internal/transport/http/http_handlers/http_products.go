package http_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ferreteria-inventario-ventas/internal/domain"
)

// Products maneja:
// GET  -> listar productos
// POST -> crear producto

// Products godoc
// @Summary Listar o crear productos
// @Description GET lista productos, POST crea producto
// @Tags Products
// @Accept json
// @Produce json
// @Param product body domain.Product false "Producto (solo POST)"
// @Success 200 {array} domain.Product
// @Success 201 {object} domain.Product
// @Router /api/products [get]
// @Router /api/products [post]
func (h *Handlers) Products(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodPut:
		idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		var input domain.Product
		json.NewDecoder(r.Body).Decode(&input)

		err := h.ProductsSvc.Update(id, &input)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, input)

	case http.MethodDelete:
		idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		err := h.ProductsSvc.Delete(id)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, map[string]string{"deleted": "ok"})

	case http.MethodGet:
		list, err := h.ProductsSvc.List()
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, list)

	case http.MethodPost:
		var input domain.Product

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "JSON inválido"})
			return
		}

		err = h.ProductsSvc.Create(&input)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 201, input)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
