package main

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
