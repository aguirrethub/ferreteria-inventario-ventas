package sqlite

import (
	"database/sql"

	"ferreteria-inventario-ventas/internal/domain"
)

// ClientRepo maneja las operaciones de base de datos para clientes.
type ClientRepo struct {
	db *sql.DB
}

// Constructor del repositorio.
func NewClientRepo(db *sql.DB) *ClientRepo {
	return &ClientRepo{db: db}
}

// Create inserta un nuevo cliente en la base de datos.
func (r *ClientRepo) Create(c *domain.Client) error {

	result, err := r.db.Exec(
		`INSERT INTO clients(nombre, cedula, email) VALUES(?,?,?)`,
		c.Nombre, c.Cedula, c.Email,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	c.ID = id

	return nil
}

// List devuelve todos los clientes.
func (r *ClientRepo) List() ([]domain.Client, error) {

	rows, err := r.db.Query(`SELECT id, nombre, cedula, email FROM clients ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []domain.Client

	for rows.Next() {
		var c domain.Client
		err := rows.Scan(&c.ID, &c.Nombre, &c.Cedula, &c.Email)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return clients, nil
}

func (r *ClientRepo) Update(id int64, c *domain.Client) error {
	_, err := r.db.Exec(
		`UPDATE clients SET nombre=?, cedula=?, email=? WHERE id=?`,
		c.Nombre, c.Cedula, c.Email, id,
	)
	return err
}

func (r *ClientRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM clients WHERE id=?`, id)
	return err
}
