package domain

// Client representa un cliente de la ferretería.
// Contiene los datos básicos necesarios para registrar ventas.
type Client struct {
	ID     int64  `json:"id"`     // Identificador único en la base de datos
	Nombre string `json:"nombre"` // Nombre completo del cliente
	Cedula string `json:"cedula"` // Número de cédula (único)
	Email  string `json:"email"`  // Correo electrónico
}
