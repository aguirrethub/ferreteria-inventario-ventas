package domain

import "time"

// SaleItem representa un producto dentro de una venta.
// Cada venta puede tener varios productos.
type SaleItem struct {
	ProductID      int64   `json:"product_id"`      // ID del producto vendido
	Cantidad       int     `json:"cantidad"`        // Cantidad vendida
	PrecioUnitario float64 `json:"precio_unitario"` // Precio al momento de la venta
	Subtotal       float64 `json:"subtotal"`        // Cantidad * PrecioUnitario
}

// Sale representa la cabecera de una venta.
type Sale struct {
	ID         int64      `json:"id"`
	ClientID   int64      `json:"client_id"`
	ClientName string     `json:"client_name"` // 👈 NUEVO
	Fecha      time.Time  `json:"fecha"`
	Total      float64    `json:"total"`
	Items      []SaleItem `json:"items"`
}
