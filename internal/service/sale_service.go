package service

import "ferreteria-inventario-ventas/internal/domain"

// Interfaz que debe cumplir el repositorio de ventas.
type SaleRepository interface {
	CreateSaleTx(clientID int64, items []domain.SaleItem) (*domain.Sale, error)
	ListSales() ([]domain.Sale, error)
	GetSaleDetail(saleID int64) (*domain.Sale, error)

	ClientExists(id int64) (bool, error)
	ProductExists(id int64) (bool, error)

	// NUEVOS MÉTODOS DE REPORTE
	VentasHoy() (int, float64, error)
	TopProductos() ([]map[string]interface{}, error)
}

// SaleService contiene la lógica de negocio para ventas.
type SaleService struct {
	repo SaleRepository
}

// Constructor del servicio.
func NewSaleService(r SaleRepository) *SaleService {
	return &SaleService{repo: r}
}

// Create valida los datos antes de registrar la venta.
func (s *SaleService) Create(clientID int64, items []domain.SaleItem) (*domain.Sale, error) {

	if clientID <= 0 || len(items) == 0 {
		return nil, domain.ErrInvalidInput
	}

	exists, err := s.repo.ClientExists(clientID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrNotFound
	}

	for _, item := range items {

		if item.ProductID <= 0 || item.Cantidad <= 0 || item.PrecioUnitario <= 0 {
			return nil, domain.ErrInvalidInput
		}

		ok, err := s.repo.ProductExists(item.ProductID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrNotFound
		}
	}

	return s.repo.CreateSaleTx(clientID, items)
}

func (s *SaleService) List() ([]domain.Sale, error) {
	return s.repo.ListSales()
}

func (s *SaleService) Detail(id int64) (*domain.Sale, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetSaleDetail(id)
}

// NUEVOS MÉTODOS DE REPORTE

func (s *SaleService) VentasHoy() (int, float64, error) {
	return s.repo.VentasHoy()
}

func (s *SaleService) TopProductos() ([]map[string]interface{}, error) {
	return s.repo.TopProductos()
}
