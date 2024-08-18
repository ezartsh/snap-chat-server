package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"snap_chat_server/models"

	"github.com/iancoleman/strcase"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Group struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type GroupRequestBody struct {
	GroupName string `json:"name"`
}

type LeaveGroupRequestBody struct {
	GroupKey string `json:"key"`
}

type JoinGroupRequestBody struct {
	GroupKey string `json:"key"`
}

func GetGroupList(db *gorm.DB, authUser *AuthSession, r *http.Request) ([]Group, error) {
	groups := []Group{}

	if err := db.Model(models.Group{}).
		Find(&groups).Error; err != nil {
		return groups, err
	}

	return groups, nil
}

func CreateGroup(db *gorm.DB, authUser *AuthSession, r *http.Request) error {
	var data GroupRequestBody

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	groupKey := strcase.ToKebab(data.GroupName)

	// Check if auth user already have this incoming new group in list.
	var existedGroup int64
	if err := db.Model(models.Group{}).Where("key = ?", groupKey).Count(&existedGroup).Error; err != nil {
		return err
	}

	if existedGroup > 0 {
		return errors.New("This group already registered.")
	}

	newRoom := models.Room{
		RoomUID:  xid.New().String(),
		RoomType: "group",
		Name:     data.GroupName,
	}

	db.Transaction(func(tx *gorm.DB) error {

		if errCreate := tx.Model(models.Room{}).Create(&newRoom).Error; errCreate != nil {
			return errCreate
		}

		if errCreate := tx.Model(models.Group{}).Create(&models.Group{
			RoomID: int(*newRoom.ID),
			Key:    groupKey,
			Name:   data.GroupName,
		}).Error; errCreate != nil {
			return errCreate
		}

		if errCreate := tx.Model(models.RoomAudience{}).Create(&models.RoomAudience{
			RoomID: *newRoom.ID,
			UserID: authUser.ID,
		}).Error; errCreate != nil {
			return errCreate
		}

		return nil
	})

	return nil
}

func LeaveGroup(db *gorm.DB, authUser *AuthSession, r *http.Request) error {
	var data LeaveGroupRequestBody

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	// Check if auth user already have this incoming new group in list.
	var group models.Group
	var room models.Room
	var roomAudience models.RoomAudience
	if err := db.Model(models.Group{}).Where("key = ?", data.GroupKey).Take(&group).Error; err != nil {
		return err
	}

	if err := db.Model(models.Room{}).Where("id = ?", group.RoomID).Take(&room).Error; err != nil {
		return err
	}

	if err := db.Model(models.RoomAudience{}).Where("room_id = ? AND user_id = ?", room.ID, authUser.ID).Take(&roomAudience).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("You're not in the group or already left the group")
		}
		return err
	}

	if err := db.Model(models.RoomAudience{}).Delete(&roomAudience, "id = ?", roomAudience.ID).Error; err != nil {
		return err
	}

	return nil
}

func JoinGroup(db *gorm.DB, authUser *AuthSession, r *http.Request) error {
	var data JoinGroupRequestBody

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	// Check if auth user already have this incoming new group in list.
	var group models.Group
	var room models.Room
	var existedAudience int64
	if err := db.Model(models.Group{}).Where("key = ?", data.GroupKey).Take(&group).Error; err != nil {
		return err
	}

	if err := db.Model(models.Room{}).Where("id = ?", group.RoomID).Take(&room).Error; err != nil {
		return err
	}

	if err := db.Model(models.RoomAudience{}).Where("room_id = ? AND user_id = ?", room.ID, authUser.ID).Count(&existedAudience).Error; err != nil {
		return err
	}

	if existedAudience == 0 {
		if errCreate := db.Model(models.RoomAudience{}).Create(&models.RoomAudience{
			RoomID: *room.ID,
			UserID: authUser.ID,
		}).Error; errCreate != nil {
			return errCreate
		}
	} else {
		return errors.New("You're already join the group.")
	}

	return nil
}
