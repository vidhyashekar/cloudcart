package handler

import (
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	fmt.Fprintln(w, "auth-service is healthy")
}
