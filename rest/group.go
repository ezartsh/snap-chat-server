package rest

import (
	"fmt"
	"net/http"
	"snap_chat_server/logger"
	"snap_chat_server/services"

	"gorm.io/gorm"
)

func GetGroup(db *gorm.DB) {

	http.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("GET", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		groups, err := services.GetGroupList(db, authSession, r)

		if err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(groups)

		jsonResponse(w, http.StatusOK, groups)
		return

	})
}

func CreateGroup(db *gorm.DB) {

	http.HandleFunc("/groups/create", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := services.CreateGroup(db, authSession, r); err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, http.StatusCreated, nil)
		return

	})
}

func JoinGroup(db *gorm.DB) {

	http.HandleFunc("/groups/join", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := services.JoinGroup(db, authSession, r); err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, http.StatusOK, nil)
		return

	})
}

func LeaveGroup(db *gorm.DB) {

	http.HandleFunc("/groups/leave", func(w http.ResponseWriter, r *http.Request) {
		if err := verifyRestMethod("POST", r); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		authSession, err := verifyShouldAuthenticated(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := services.LeaveGroup(db, authSession, r); err != nil {
			logger.AppLog.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, http.StatusOK, nil)
		return

	})
}
