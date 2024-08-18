package services

import (
	"errors"
	"snap_chat_server/models"

	"gorm.io/gorm"
)

type GroupChatSender struct {
	Name     string
	Username string
}

type GroupChat struct {
	RoomID   string
	RoomType string
	Sender   GroupChatSender
	Target   string
}

func NewGroupChat(sender GroupChatSender, target string) *GroupChat {
	return &GroupChat{
		RoomType: "private",
		Sender:   sender,
		Target:   target,
	}
}

func (c *GroupChat) Find(db *gorm.DB, authSession AuthSession) (models.Room, []string, error) {
	var room models.Room
	var groupTarget models.Group
	var targetAudiences = []string{}

	if err := db.Model(models.Group{}).
		Where("key = ?", c.Target).
		Take(&groupTarget).Error; err != nil {
		return models.Room{}, targetAudiences, err
	}

	if err := db.Model(models.Room{}).
		Where("id = ?", groupTarget.RoomID).
		Take(&room).Error; err != nil {

		return models.Room{}, targetAudiences, err
	}

	var isAuthUserJoinedGroup int64

	if err := db.Model(models.RoomAudience{}).
		Where("room_id = ?", *room.ID).
		Where("user_id = ?", authSession.ID).
		Count(&isAuthUserJoinedGroup).Error; err != nil {

		return models.Room{}, targetAudiences, err
	}

	if isAuthUserJoinedGroup == 0 {
		return models.Room{}, targetAudiences, errors.New("You're not join to this group yet.")
	}

	if err := db.Model(models.RoomAudience{}).
		Select("users.username").
		Where("room_id = ?", *room.ID).
		Where("user_id <> ?", authSession.ID).
		Joins("JOIN rooms ON room_id = rooms.id AND rooms.room_type = 'group'").
		Joins("JOIN users ON user_id = users.id").
		Find(&targetAudiences).Error; err != nil {

		return models.Room{}, targetAudiences, err
	}

	return room, targetAudiences, nil
}
