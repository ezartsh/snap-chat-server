package rest

import (
	"encoding/json"
	"net/http"
)

func HealthCheck() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("ok")
		return
	})
}
