package msg

import (
	"context"

	"ar-app-api/dal/data"
	"ar-app-api/model"
	"ar-app-api/util"
)

func SaveMessage(ctx context.Context, msg *model.Message) (*model.Message, error) {
	// todo 存
	message, err := data.CreateMessage(ctx, util.ToMsgSchema(msg))
	if err != nil {
		return nil, err
	}
	return message.ToType(), nil
}

func SendFinish(ctx context.Context, msg *model.Message) (*model.Message, error) {
	msg.IsSend = true
	message, err := data.UpdateMessage(ctx, util.ToMsgSchema(msg))
	if err != nil {
		return nil, err
	}
	return message.ToType(), nil
}
