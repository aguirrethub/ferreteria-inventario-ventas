package sqlite

import (
	"database/sql"
	"time"

	"ferreteria-inventario-ventas/internal/domain"
)

// SaleRepo maneja las operaciones relacionadas a ventas.
type SaleRepo struct {
	db *sql.DB
}

// Constructor del repositorio.
func NewSaleRepo(db *sql.DB) *SaleRepo {
	return &SaleRepo{db: db}
}

// CreateSaleTx crea una venta completa usando transacción.
// 1) Inserta la cabecera
// 2) Inserta los productos vendidos
// 3) Descuenta el stock
func (r *SaleRepo) CreateSaleTx(clientID int64, items []domain.SaleItem) (*domain.Sale, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	fecha := time.Now()

	var total float64

	// Calcular total y subtotales
	for i := range items {
		items[i].Subtotal = float64(items[i].Cantidad) * items[i].PrecioUnitario
		total += items[i].Subtotal
	}

	// Insertar cabecera de venta
	result, err := tx.Exec(
		`INSERT INTO sales(client_id, fecha, total) VALUES(?,?,?)`,
		clientID,
		fecha.Format(time.RFC3339),
		total,
	)
	if err != nil {
		return nil, err
	}

	saleID, _ := result.LastInsertId()

	// Insertar detalle y descontar stock
	for _, item := range items {

		// Descuento de stock validando que haya suficiente
		res, err := tx.Exec(
			`UPDATE products 
			 SET stock = stock - ? 
			 WHERE id = ? AND stock >= ?`,
			item.Cantidad,
			item.ProductID,
			item.Cantidad,
		)
		if err != nil {
			return nil, err
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return nil, domain.ErrInsufficientStock
		}

		// Insertar detalle
		_, err = tx.Exec(
			`INSERT INTO sale_items(sale_id, product_id, cantidad, precio_unitario, subtotal)
			 VALUES(?,?,?,?,?)`,
			saleID,
			item.ProductID,
			item.Cantidad,
			item.PrecioUnitario,
			item.Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &domain.Sale{
		ID:       saleID,
		ClientID: clientID,
		Fecha:    fecha,
		Total:    total,
		Items:    items,
	}, nil
}

// ListSales devuelve todas las ventas registradas.
// ListSales devuelve todas las ventas registradas.
func (r *SaleRepo) ListSales() ([]domain.Sale, error) {

	rows, err := r.db.Query(`
		SELECT s.id, s.client_id, c.nombre, s.fecha, s.total
		FROM sales s
		JOIN clients c ON c.id = s.client_id
		ORDER BY s.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []domain.Sale

	for rows.Next() {
		var s domain.Sale
		var fechaStr string

		if err := rows.Scan(&s.ID, &s.ClientID, &s.ClientName, &fechaStr, &s.Total); err != nil {
			return nil, err
		}

		if t, e := time.Parse(time.RFC3339, fechaStr); e == nil {
			s.Fecha = t
		}

		sales = append(sales, s)
	}

	return sales, rows.Err()
}

// GetSaleDetail devuelve una venta con sus items.
func (r *SaleRepo) GetSaleDetail(saleID int64) (*domain.Sale, error) {

	// 1) Cabecera
	var s domain.Sale
	var fechaStr string

	err := r.db.QueryRow(
		`SELECT id, client_id, fecha, total FROM sales WHERE id = ?`,
		saleID,
	).Scan(&s.ID, &s.ClientID, &fechaStr, &s.Total)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Intentar parsear fecha (si falla, no rompe)
	if t, e := time.Parse(time.RFC3339, fechaStr); e == nil {
		s.Fecha = t
	}

	// 2) Items
	rows, err := r.db.Query(
		`SELECT product_id, cantidad, precio_unitario, subtotal
		 FROM sale_items
		 WHERE sale_id = ?
		 ORDER BY id ASC`,
		saleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it domain.SaleItem
		if err := rows.Scan(&it.ProductID, &it.Cantidad, &it.PrecioUnitario, &it.Subtotal); err != nil {
			return nil, err
		}
		s.Items = append(s.Items, it)
	}

	return &s, rows.Err()
}

func (r *SaleRepo) ClientExists(id int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM clients WHERE id = ?`, id).Scan(&count)
	return count > 0, err
}

func (r *SaleRepo) ProductExists(id int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM products WHERE id = ?`, id).Scan(&count)
	return count > 0, err
}

// VentasHoy devuelve total de ventas y monto del día actual.
func (r *SaleRepo) VentasHoy() (int, float64, error) {

	rows, err := r.db.Query(`
		SELECT COUNT(*), IFNULL(SUM(total),0)
		FROM sales
		WHERE DATE(fecha) = DATE('now')
	`)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var count int
	var total float64

	if rows.Next() {
		if err := rows.Scan(&count, &total); err != nil {
			return 0, 0, err
		}
	}

	return count, total, nil
}

// TopProductos devuelve productos más vendidos.
func (r *SaleRepo) TopProductos() ([]map[string]interface{}, error) {

	rows, err := r.db.Query(`
		SELECT p.nombre, SUM(si.cantidad) as total_vendido
		FROM sale_items si
		JOIN products p ON p.id = si.product_id
		GROUP BY p.nombre
		ORDER BY total_vendido DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var nombre string
		var total int
		if err := rows.Scan(&nombre, &total); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"producto": nombre,
			"cantidad": total,
		})
	}

	return result, nil
}
