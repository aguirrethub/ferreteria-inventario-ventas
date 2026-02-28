Ferretería Inventario Ventas

Sistema web para gestión de inventario y ventas orientado a un negocio tipo ferretería.
Incluye API REST en Go, persistencia en SQLite, transacciones para ventas (descuento de stock atómico), reportes y una interfaz web para operar productos y ventas.

✨ Funcionalidades principales
Inventario (Productos)

Crear productos (nombre, precio, stock)

Listar productos

Editar productos

Eliminar productos

Validación de datos y respuestas JSON

Ventas

Crear venta con múltiples items

Registro de cabecera + detalle de items

Descuento de stock dentro de transacción

Consulta de ventas (cabecera)

Consulta de detalle de venta por ID

Reportes

Ventas del día

Top productos (más vendidos / más demandados)

Documentación

Swagger operativo para probar servicios desde el navegador

🧠 Cómo funciona (visión técnica)

Este proyecto está hecho con un enfoque por capas, para que el código sea mantenible y profesional:

Domain: modelos del negocio (Producto, Venta, Items, etc.)

Service / Use cases: reglas del negocio (validaciones, lógica de venta)

Storage (SQLite): acceso a datos con database/sql

Transport (HTTP): API REST, handlers y routing

Web UI: HTML/CSS/JS consumiendo la API

Punto crítico: Ventas con transacción

Cuando confirmas una venta, el sistema hace esto en una sola transacción:

Inserta la cabecera de venta

Inserta los items vendidos (detalle)

Descuenta stock por cada producto con operación segura

Si algo falla (sin stock, ID inválido, error SQL) => rollback (no se guarda nada a medias)

Esto evita ventas “fantasma” y mantiene la BD consistente.

🧱 Tecnologías usadas

Go (backend)

SQLite (base de datos)

Swagger / OpenAPI (documentación y pruebas)

HTML + CSS + JavaScript (interfaz web)

database/sql para consultas y transacciones

📁 Estructura del proyecto (resumen)

Los nombres pueden variar según tu repo, pero el concepto es este:

cmd/ → punto de entrada (arranque del servidor)

internal/domain/ → entidades del negocio

internal/service/ → reglas, validaciones, casos de uso

internal/storage/sqlite/ → repositorios SQLite (productos, ventas, reportes)

internal/transport/http/ → handlers, rutas, middleware básico (si aplica)

web/ o ui/ → interfaz HTML/CSS/JS

db/ → archivo .db o scripts de inicialización

⚙️ Requisitos

Go 1.20+ (recomendado)

SQLite (normalmente ya viene integrado si usas archivo .db)

Navegador web (para UI y Swagger)

▶️ Cómo ejecutar el proyecto
1) Clonar el repositorio
git clone <URL_DE_TU_REPO>
cd <NOMBRE_DEL_PROYECTO>
2) Instalar dependencias (si aplica)
go mod tidy
3) Ejecutar el servidor
go run ./cmd/api

Al iniciar, el sistema:

levanta el servidor HTTP

prepara la base SQLite (tablas si no existen)

expone API + UI web

🌐 Rutas principales
UI Web

/ o /ui → pantalla principal (productos / ventas)

/sales o /ui/sales.html → módulo de ventas (según tu estructura)

Swagger

/swagger/index.html (o ruta equivalente) → documentación interactiva

🔌 API REST (Servicios Web)

El proyecto cumple el requisito académico de 8+ servicios web con serialización JSON.

Productos

GET /api/products → listar productos

POST /api/products → crear producto

GET /api/products/{id} → obtener producto por ID

PUT /api/products/{id} → actualizar producto

DELETE /api/products/{id} → eliminar producto

Ventas

GET /api/sales → listar ventas (cabecera)

POST /api/sales → crear venta (transacción: cabecera + items + descuento stock)

GET /api/sales/{id} → detalle de venta (cabecera + items)

Reportes

GET /api/report/ventas-hoy → total ventas del día + resumen

GET /api/report/top-productos → productos más vendidos

(Si tu proyecto tiene nombres exactos distintos, cambia únicamente las rutas, pero el README ya está listo.)

✅ Ejemplo de venta (JSON)
Crear venta

POST /api/sales

{
  "cliente_nombre": "Juan Pérez",
  "items": [
    { "product_id": 1, "qty": 2 },
    { "product_id": 3, "qty": 1 }
  ]
}

Resultado esperado:

se crea la venta

se guarda el detalle

se descuenta stock

responde con JSON de confirmación (ID, totales, fecha, etc.)

🧪 Pruebas rápidas (recomendado)

Probar endpoints desde Swagger

Crear 2–3 productos

Hacer una venta con 2 items

Verificar:

que el stock baja

que la venta aparece en el listado

que el detalle de venta muestra items y totales

que el reporte de ventas del día refleja el movimiento

🎓 Justificación del proyecto (enfoque académico)

Este sistema se eligió porque simula un escenario real de negocio (ferretería), donde se requieren:

control de inventario (stock, productos)

registro formal de ventas (cabecera + detalle)

consistencia de datos mediante transacciones

exposición de funcionalidades como Servicios Web REST

serialización JSON

reportes operativos para toma de decisiones

👤 Autor

Gabriel Aguirre Román
Ingeniería en Software – UIDE