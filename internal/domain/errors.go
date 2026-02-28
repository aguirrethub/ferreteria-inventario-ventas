package domain

import "errors"

// Errores estándar del sistema.
// Nos permiten manejar validaciones de forma clara.
var (
	ErrNotFound          = errors.New("not found")          // No existe el registro
	ErrInvalidInput      = errors.New("invalid input")      // Datos incorrectos
	ErrConflict          = errors.New("conflict")           // Conflicto (ej: cédula repetida)
	ErrInsufficientStock = errors.New("insufficient stock") // Stock insuficiente
)
