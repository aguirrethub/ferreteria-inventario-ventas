package service

import "ferreteria-inventario-ventas/internal/domain"

// Interfaz que define lo que el repositorio debe implementar.
type ClientRepository interface {
	Create(*domain.Client) error
	List() ([]domain.Client, error)
	Update(id int64, c *domain.Client) error
	Delete(id int64) error
}

// ClientService contiene la lógica de negocio para clientes.
type ClientService struct {
	repo ClientRepository
}

// Constructor del servicio.
func NewClientService(r ClientRepository) *ClientService {
	return &ClientService{repo: r}
}

// Create valida los datos antes de guardar.
func (s *ClientService) Create(c *domain.Client) error {

	if c.Nombre == "" || c.Cedula == "" || c.Email == "" {
		return domain.ErrInvalidInput
	}

	return s.repo.Create(c)
}

// List devuelve todos los clientes.
func (s *ClientService) List() ([]domain.Client, error) {
	return s.repo.List()
}

func (s *ClientService) Update(id int64, c *domain.Client) error {
	if id <= 0 || c.Nombre == "" || c.Cedula == "" || c.Email == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.Update(id, c)
}

func (s *ClientService) Delete(id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(id)
}
