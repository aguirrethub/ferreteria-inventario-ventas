package sqlite

import (
	"database/sql"
	"os"
)

// Migrate ejecuta el archivo schema.sql.
// Su función es crear las tablas si no existen.
func Migrate(db *sql.DB, schemaPath string) error {

	// Leer el archivo SQL
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	// Ejecutar el SQL en la base de datos
	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	return nil
}
