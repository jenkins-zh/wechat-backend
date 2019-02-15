package health

import (
	"net/http"
)

// SimpleHealthHandler indicate server status
func SimpleHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
