package data

import (
	"context"

	"ar-app-api/dal"
	"ar-app-api/dal/schema"
)

func CreateMessage(ctx context.Context, msg *schema.Message) (*schema.Message, error) {
	err := dal.DB.WithContext(ctx).Create(msg).Error
	return msg, err
}

func UpdateMessage(ctx context.Context, msg *schema.Message) (*schema.Message, error) {
	err := dal.DB.WithContext(ctx).Save(msg).Error
	return msg, err
}
