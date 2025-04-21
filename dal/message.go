package dal

import (
	"context"

	"learning-assistant/dal/schema"
)

func CreateMessage(ctx context.Context, msg *schema.Message) (*schema.Message, error) {
	err := DB.WithContext(ctx).Create(msg).Error
	return msg, err
}

func UpdateMessage(ctx context.Context, msg *schema.Message) (*schema.Message, error) {
	err := DB.WithContext(ctx).Save(msg).Error
	return msg, err
}
