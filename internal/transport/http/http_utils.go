package httptransport

import (
	"encoding/json"
	"net/http"
)

// WriteJSON envía una respuesta JSON al cliente.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// ReadJSON lee el body de la petición y lo convierte en struct.
func ReadJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
