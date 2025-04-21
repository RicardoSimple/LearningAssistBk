package util

import (
	"gorm.io/gorm"
	"learning-assistant/consts"

	"learning-assistant/dal/schema"
	"learning-assistant/model"
)

func ToUserSchema(u *model.User) *schema.User {
	return &schema.User{
		Model: gorm.Model{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Username:    u.Username,
		Password:    u.Password,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Gender:      u.Gender,
		DateOfBirth: u.DateOfBirth,
		Address:     u.Address,
		City:        u.City,
		State:       u.State,
		Country:     u.Country,
		PostalCode:  u.PostalCode,
		Status:      u.Status,
		LastLogin:   u.LastLogin,
		ClassStage:  u.ClassStage,
		UserType:    consts.UserTypeToIntMap[u.UserType],
	}
}

func ToMsgSchema(msg *model.Message) *schema.Message {
	return &schema.Message{
		Model: gorm.Model{
			ID:        msg.ID,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		},
		SenderID:    msg.SenderID,
		ReceiverID:  msg.ReceiverID,
		RoomID:      msg.RoomID,
		Content:     msg.Content,
		MessageType: msg.MessageType,
		Timestamp:   msg.Timestamp,
		IsRead:      msg.IsRead,
		IsSend:      msg.IsSend,
	}
}
