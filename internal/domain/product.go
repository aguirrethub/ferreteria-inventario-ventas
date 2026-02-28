package domain

// Product representa un producto del inventario de la ferretería.
// Ejemplo: "Saco de cemento 50kg", stock 300, precio 8.50
type Product struct {
	ID     int64   `json:"id"`     // Identificador único en la base de datos
	Nombre string  `json:"nombre"` // Nombre del producto
	Stock  int     `json:"stock"`  // Cantidad disponible en inventario
	Precio float64 `json:"precio"` // Precio unitario del producto
}
