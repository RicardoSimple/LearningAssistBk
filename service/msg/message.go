package msg

import (
	"context"
	"learning-assistant/dal"

	"learning-assistant/model"
	"learning-assistant/util"
)

func SaveMessage(ctx context.Context, msg *model.Message) (*model.Message, error) {
	// todo 存
	message, err := dal.CreateMessage(ctx, util.ToMsgSchema(msg))
	if err != nil {
		return nil, err
	}
	return message.ToType(), nil
}

func SendFinish(ctx context.Context, msg *model.Message) (*model.Message, error) {
	msg.IsSend = true
	message, err := dal.UpdateMessage(ctx, util.ToMsgSchema(msg))
	if err != nil {
		return nil, err
	}
	return message.ToType(), nil
}
