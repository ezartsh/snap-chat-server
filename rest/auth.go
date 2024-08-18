package rest

import (
	"encoding/json"
	"net/http"
	"snap_chat_server/services"

	"gorm.io/gorm"
)

func AccountRegister(db *gorm.DB) {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		user, err := services.Register(db, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
		return
	})
}

func AccountLogin(db *gorm.DB) {

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		resp, err := services.Login(db, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return

	})

}
