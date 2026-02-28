package service

import "ferreteria-inventario-ventas/internal/domain"

// Interfaz que debe cumplir el repositorio de productos.
type ProductRepository interface {
	Create(*domain.Product) error
	List() ([]domain.Product, error)
	Update(id int64, p *domain.Product) error
	Delete(id int64) error
}

// ProductService contiene la lógica de negocio para productos.
type ProductService struct {
	repo ProductRepository
}

// Constructor del servicio.
func NewProductService(r ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

// Create valida datos antes de guardar.
func (s *ProductService) Create(p *domain.Product) error {

	if p.Nombre == "" || p.Stock < 0 || p.Precio <= 0 {
		return domain.ErrInvalidInput
	}

	return s.repo.Create(p)
}

// List devuelve todos los productos.
func (s *ProductService) List() ([]domain.Product, error) {
	return s.repo.List()
}

func (s *ProductService) Update(id int64, p *domain.Product) error {
	if id <= 0 || p.Nombre == "" || p.Stock < 0 || p.Precio <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.Update(id, p)
}

func (s *ProductService) Delete(id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(id)
}
