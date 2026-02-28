package http_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ferreteria-inventario-ventas/internal/domain"
)

// Sales godoc
// @Summary Listar o crear ventas
// @Description GET lista ventas, POST crea venta
// @Tags Sales
// @Accept json
// @Produce json
// @Param sale body domain.Sale false "Venta (solo POST)"
// @Success 200 {array} domain.Sale
// @Success 201 {object} domain.Sale
// @Router /api/sales [get]
// @Router /api/sales [post]
func (h *Handlers) Sales(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		list, err := h.SalesSvc.List()
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, list)

	case http.MethodPost:
		var input struct {
			ClientID int64             `json:"client_id"`
			Items    []domain.SaleItem `json:"items"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "JSON inválido"})
			return
		}

		sale, err := h.SalesSvc.Create(input.ClientID, input.Items)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 201, sale)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SaleDetail godoc
// @Summary Obtener detalle de venta
// @Description Devuelve una venta con sus items
// @Tags Sales
// @Produce json
// @Param id path int true "ID de la venta"
// @Success 200 {object} domain.Sale
// @Router /api/sales/{id} [get]
func (h *Handlers) SaleDetail(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/sales/")
	if idStr == "" {
		writeJSON(w, 400, map[string]string{"error": "id requerido"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, 400, map[string]string{"error": "id inválido"})
		return
	}

	sale, err := h.SalesSvc.Detail(id)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			writeJSON(w, 404, map[string]string{"error": "venta no encontrada"})
		case domain.ErrInvalidInput:
			writeJSON(w, 400, map[string]string{"error": "id inválido"})
		default:
			writeJSON(w, 500, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, 200, sale)
}

// ReportVentasHoy godoc
// @Summary Ventas del día
// @Description Devuelve cantidad de ventas y total del día
// @Tags Report
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/report/ventas-hoy [get]
func (h *Handlers) ReportVentasHoy(w http.ResponseWriter, r *http.Request) {

	count, total, err := h.SalesSvc.VentasHoy()
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"ventas": count,
		"total":  total,
	})
}

// ReportTopProductos godoc
// @Summary Top productos vendidos
// @Description Devuelve los 5 productos más vendidos
// @Tags Report
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/report/top-productos [get]
func (h *Handlers) ReportTopProductos(w http.ResponseWriter, r *http.Request) {

	data, err := h.SalesSvc.TopProductos()
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, data)
}
