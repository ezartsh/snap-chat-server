package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"snap_chat_server/models"

	"gorm.io/gorm"
)

type Contact struct {
	Username string `json:"username"`
}

type ContactRequestBody struct {
	Username string `json:"username"`
}

func GetContactList(db *gorm.DB, authUser *AuthSession, r *http.Request) ([]Contact, error) {
	contacts := []Contact{}

	if err := db.Model(models.Contact{}).
		Select("users.username").
		Where("user_id = ?", authUser.ID).
		Joins("JOIN users ON contact_id = users.id").
		Find(&contacts).Error; err != nil {
		return contacts, err
	}

	return contacts, nil
}

func AddContact(db *gorm.DB, authUser *AuthSession, r *http.Request) error {
	var userContact models.User
	var data ContactRequestBody

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	if err := db.Model(models.User{}).Where("username = ?", data.Username).Take(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("Cannot find user with name %s", data.Username))
		}
		return err
	}

	if authUser.Username == userContact.Username {
		return errors.New("Cannot add yourself to your contact.")
	}

	// Check if auth user already have this incoming new contact in list.
	var existedContact int64
	if err := db.Model(models.Contact{}).Where("user_id = ? AND contact_id = ?", authUser.ID, userContact.ID).Count(&existedContact).Error; err != nil {
		return err
	}

	if existedContact > 0 {
		return errors.New("This user already registered in your contact list.")
	}

	newContact := models.Contact{
		UserID:    authUser.ID,
		ContactID: *userContact.ID,
	}

	if err := db.Model(models.Contact{}).Create(&newContact).Error; err != nil {
		return err
	}

	return nil
}

func RemoveContact(db *gorm.DB, authUser *AuthSession, r *http.Request) error {
	var userContact models.User
	var data ContactRequestBody

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	if err := db.Model(models.User{}).Where("username = ?", data.Username).Take(&userContact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("Cannot find user with name %s", data.Username))
		}
		return err
	}

	// Check if auth user already have this incoming new contact in list.
	var existedContact int64
	if err := db.Model(models.Contact{}).Where("user_id = ? AND contact_id = ?", authUser.ID, userContact.ID).Count(&existedContact).Error; err != nil {
		return err
	}

	if existedContact == 0 {
		return errors.New("This user not in your registered contact list.")
	}

	deletedContact := models.Contact{
		UserID:    authUser.ID,
		ContactID: *userContact.ID,
	}

	if err := db.Model(models.Contact{}).Where("user_id = ? AND contact_id = ?", authUser.ID, userContact.ID).Delete(&deletedContact).Error; err != nil {
		return err
	}

	return nil
}
