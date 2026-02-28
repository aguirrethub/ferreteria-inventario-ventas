-- ================================
-- TABLA CLIENTES
-- ================================
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    cedula TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL
);

-- ================================
-- TABLA PRODUCTOS
-- ================================
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL UNIQUE,
    stock INTEGER NOT NULL,
    precio REAL NOT NULL
);

-- ================================
-- TABLA VENTAS (CABECERA)
-- ================================
CREATE TABLE IF NOT EXISTS sales (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    fecha TEXT NOT NULL,
    total REAL NOT NULL,
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

-- ================================
-- TABLA DETALLE DE VENTA
-- ================================
CREATE TABLE IF NOT EXISTS sale_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sale_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    cantidad INTEGER NOT NULL,
    precio_unitario REAL NOT NULL,
    subtotal REAL NOT NULL,
    FOREIGN KEY (sale_id) REFERENCES sales(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);