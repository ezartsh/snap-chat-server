package services

import (
	"errors"
	"snap_chat_server/models"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type PrivateChatSender struct {
	Name     string
	Username string
}

type PrivateChat struct {
	RoomID   string
	RoomType string
	Sender   PrivateChatSender
	Target   string
}

func NewPrivateChat(sender PrivateChatSender, target string) *PrivateChat {
	return &PrivateChat{
		RoomType: "private",
		Sender:   sender,
		Target:   target,
	}
}

func (c *PrivateChat) FindOrCreate(db *gorm.DB, authSession AuthSession) (models.Room, error) {
	var roomAudience models.RoomAudience
	var room models.Room
	var userTarget models.User

	if err := db.Model(models.User{}).
		Where("username = ?", c.Target).
		Take(&userTarget).Error; err != nil {
		return models.Room{}, err
	}

	if err := db.Model(models.RoomAudience{}).
		Where("user_id IN ?", []int{int(authSession.ID), int(*userTarget.ID)}).
		Joins("JOIN rooms ON room_id = rooms.id AND rooms.room_type = 'private'").
		Take(&roomAudience).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			newRoom := models.Room{
				RoomUID:  xid.New().String(),
				RoomType: "private",
				Name:     "Private",
			}

			db.Transaction(func(tx *gorm.DB) error {

				if errCreate := tx.Model(models.Room{}).Create(&newRoom).Error; errCreate != nil {
					return errCreate
				}

				if errCreate := tx.Model(models.RoomAudience{}).Create(&models.RoomAudience{
					RoomID: *newRoom.ID,
					UserID: authSession.ID,
				}).Error; errCreate != nil {
					return errCreate
				}

				if errCreate := tx.Model(models.RoomAudience{}).Create(&models.RoomAudience{
					RoomID: *newRoom.ID,
					UserID: *userTarget.ID,
				}).Error; errCreate != nil {
					return errCreate
				}

				return nil
			})

			return newRoom, nil

		} else {

			return models.Room{}, err

		}
	}

	if err := db.Model(models.Room{}).
		Where("id = ?", roomAudience.RoomID).
		Take(&room).Error; err != nil {

		return models.Room{}, err
	}

	return room, nil
}
