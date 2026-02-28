package httptransport

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "ferreteria-inventario-ventas/docs"
	"ferreteria-inventario-ventas/internal/transport/http/http_handlers"
)

// NewRouter crea el router principal del sistema.
func NewRouter(h *http_handlers.Handlers) http.Handler {
	mux := http.NewServeMux()

	// =========================
	// API
	// =========================

	// Ruta para verificar que el servidor está funcionando
	mux.HandleFunc("/api/health", h.Health)

	// Clientes
	mux.HandleFunc("/api/clients", h.Clients)

	// Productos
	mux.HandleFunc("/api/products", h.Products)
	// ✅ IMPORTANTE: habilita /api/products/{id} para PUT/DELETE
	mux.HandleFunc("/api/products/", h.Products)

	// Ventas
	mux.HandleFunc("/api/sales", h.Sales)

	// Detalle de venta por ID
	mux.HandleFunc("/api/sales/", h.SaleDetail)

	// Reportes
	mux.HandleFunc("/api/report/ventas-hoy", h.ReportVentasHoy)
	mux.HandleFunc("/api/report/top-productos", h.ReportTopProductos)

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// =========================
	// UI (archivos estáticos)
	// =========================
	fs := http.FileServer(http.Dir("./web"))

	// Página raíz: redirige a Productos
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Si piden "/" => manda a products.html
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/pages/products.html", http.StatusFound)
			return
		}
		// Caso general: sirve archivos dentro de /web
		fs.ServeHTTP(w, r)
	})

	return mux
}
