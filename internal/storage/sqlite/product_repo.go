package sqlite

import (
	"database/sql"

	"ferreteria-inventario-ventas/internal/domain"
)

// ProductRepo maneja las operaciones de base de datos para productos.
type ProductRepo struct {
	db *sql.DB
}

// Constructor del repositorio.
func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create inserta un nuevo producto en la base de datos.
func (r *ProductRepo) Create(p *domain.Product) error {

	result, err := r.db.Exec(
		`INSERT INTO products(nombre, stock, precio) VALUES(?,?,?)`,
		p.Nombre, p.Stock, p.Precio,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id

	return nil
}

// List devuelve todos los productos.
func (r *ProductRepo) List() ([]domain.Product, error) {

	rows, err := r.db.Query(`SELECT id, nombre, stock, precio FROM products ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Nombre, &p.Stock, &p.Precio)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepo) Update(id int64, p *domain.Product) error {
	_, err := r.db.Exec(
		`UPDATE products SET nombre=?, stock=?, precio=? WHERE id=?`,
		p.Nombre, p.Stock, p.Precio, id,
	)
	return err
}

func (r *ProductRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id=?`, id)
	return err
}
