package write

import (
	"encoding/json"
	"net/http"
)

// Success writes a generic { success: true } response
func Success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&map[string]bool{"success": true})
}
