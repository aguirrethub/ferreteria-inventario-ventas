package main

// @title Ferretería Inventario API
// @version 1.0
// @description API REST para sistema de inventario y ventas de ferretería.
// @host localhost:8080
// @BasePath /api

import (
	"log"
	"net/http"

	"ferreteria-inventario-ventas/internal/service"
	"ferreteria-inventario-ventas/internal/storage/sqlite"
	httptransport "ferreteria-inventario-ventas/internal/transport/http"
	"ferreteria-inventario-ventas/internal/transport/http/http_handlers"
)

func main() {

	// 1️⃣ Abrir base de datos SQLite (archivo data.db)
	db, err := sqlite.OpenDB("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2️⃣ Ejecutar migraciones (crear tablas si no existen)
	if err := sqlite.Migrate(db, "migrations/schema.sql"); err != nil {
		log.Fatal(err)
	}

	// 3️⃣ Crear repositorios
	clientRepo := sqlite.NewClientRepo(db)
	productRepo := sqlite.NewProductRepo(db)
	saleRepo := sqlite.NewSaleRepo(db)

	// 4️⃣ Crear servicios (lógica de negocio)
	clientService := service.NewClientService(clientRepo)
	productService := service.NewProductService(productRepo)
	saleService := service.NewSaleService(saleRepo)

	// 5️⃣ Crear handlers HTTP
	h := &http_handlers.Handlers{
		ClientsSvc:  clientService,
		ProductsSvc: productService,
		SalesSvc:    saleService,
	}

	// 6️⃣ Crear router
	router := httptransport.NewRouter(h)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
