package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"snap_chat_server/config"
	"snap_chat_server/logger"
	"snap_chat_server/services"
	"strings"

	"gorm.io/gorm"
)

func verifyRestMethod(method string, r *http.Request) error {
	if r.Method != method {
		logger.AppLog.Error(errors.New("cannot access endpoint, method not allowed"), "Endpoint api access not allowed")
		return errors.New("cannot access endpoint, method not allowed")
	}
	return nil
}

func verifyShouldAuthenticated(r *http.Request, db *gorm.DB) (*services.AuthSession, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		logger.AppLog.Error(errors.New("No token provieded"), "Verification token failed")
		return nil, errors.New("Verification token failed")
	}

	splitToken := strings.Replace(reqToken, "Bearer ", "", 1)
	reqToken = splitToken

	claims, err := services.VerifyToken(reqToken, config.Env.SecretKey)

	if err != nil {
		logger.AppLog.Error(err, "Verification token failed")
		return nil, errors.New("Verification token failed")
	}

	username := claims["username"].(string)

	authSession, err := services.GetAuthUser(db, username, r)

	if err != nil {
		logger.AppLog.Error(err, "Get user session failed")
		return nil, err
	}

	return authSession, nil
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
