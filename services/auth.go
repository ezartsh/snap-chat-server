package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"snap_chat_server/config"
	"snap_chat_server/models"
	"snap_chat_server/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthSession struct {
	ID       uint
	Name     string
	Username string
}

type RegisterData struct {
	Name     string
	Username string
	Password string
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	Username    string `json:"username"`
}

type LoginData struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	Username    string `json:"username"`
}

func Register(db *gorm.DB, r *http.Request) (RegisterResponse, error) {
	var count int64
	var data RegisterData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return RegisterResponse{}, err
	}

	if err := db.Model(models.User{}).Where("username = ?", data.Username).Count(&count).Error; err != nil {
		return RegisterResponse{}, err
	}

	if count > 0 {
		return RegisterResponse{}, errors.New("username already registered.")
	}

	password, _ := utils.HashPassword(data.Password)

	newUser := models.User{
		Name:     data.Name,
		Username: data.Username,
		Password: password,
	}

	if err := db.Model(models.User{}).Create(&newUser).Error; err != nil {
		return RegisterResponse{}, err
	}

	newToken, err := CreateToken(newUser.Username, config.Env.SecretKey)

	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		AccessToken: newToken,
		Name:        newUser.Name,
		Username:    newUser.Username,
	}, nil
}

func Login(db *gorm.DB, r *http.Request) (LoginResponse, error) {
	var user models.User
	var data LoginData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return LoginResponse{}, err
	}

	if err := db.Model(models.User{}).Where("username = ?", data.Username).Take(&user).Error; err != nil {
		return LoginResponse{}, err
	}

	if !utils.CheckPasswordHash(data.Password, user.Password) {
		return LoginResponse{}, errors.New("credentials not matched.")
	}

	newToken, err := CreateToken(user.Username, config.Env.SecretKey)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: newToken,
		Name:        user.Name,
		Username:    user.Username,
	}, nil
}

func GetAuthUser(db *gorm.DB, username string, r *http.Request) (*AuthSession, error) {
	var authUser *AuthSession

	if err := db.Model(models.User{}).Where("username = ?", username).Take(&authUser).Error; err != nil {
		return nil, err
	}

	if authUser == nil {
		return nil, errors.New("user not found.")
	}

	return authUser, nil
}

func CreateToken(username string, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
