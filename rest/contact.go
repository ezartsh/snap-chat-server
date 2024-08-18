package rest

import (
	"net/http"
	"snap_chat_server/logger"
	"snap_chat_server/services"

	"gorm.io/gorm"
)

func GetContact(db *gorm.DB) {

	http.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("GET", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		contacts, err := services.GetContactList(db, authSession, r)

		if err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, http.StatusOK, contacts)
		return

	})
}

func AddContact(db *gorm.DB) {

	http.HandleFunc("/contacts/create", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := services.AddContact(db, authSession, r); err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, http.StatusCreated, nil)
		return

	})
}

func RemoveContact(db *gorm.DB) {

	http.HandleFunc("/contacts/remove", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := services.RemoveContact(db, authSession, r); err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, http.StatusOK, nil)
		return

	})
}
