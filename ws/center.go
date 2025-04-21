package ws

import (
	"context"
	"encoding/json"
	"strconv"

	"learning-assistant/model"
	msgService "learning-assistant/service/msg"
	"learning-assistant/util/log"
)

const (
	RetriesTimes = 10 // 重试次数
)

var messageCenter *MessageCenter

type MessageCenter struct {
}

func init() {
	messageCenter = newMessageCenter()
}
func newMessageCenter() *MessageCenter {
	return &MessageCenter{}
}

func (mc *MessageCenter) MessageControl(ctx context.Context, message string) {
	msg := &model.Message{}
	err := json.Unmarshal([]byte(message), msg)
	if err != nil {
		log.Info("[WEBSOCKET] message unmarshal fail%v", err)
	}
	tosend := mc.JudgeMessage(ctx, msg)
	// 持久化
	if !tosend {
		return
	}
	if !mc.ReceiverOnline(ctx, msg.ReceiverID) {
		mc.SendMessageLater(ctx, msg)
		return
	}
	for i := 0; i < RetriesTimes; i++ {
		err := SendMessageToUser(strconv.Itoa(int(msg.ReceiverID)), message)
		if err != nil {
			log.Info("[WEBSOCKET] send message fail %v,重试%d", err, i)
		} else {
			// 消息已发送
			go msgService.SendFinish(ctx, msg)
			return
		}
	}
	mc.SendMessageLater(ctx, msg)
	return
}

// JudgeMessage 判断消息是否该发送
func (mc *MessageCenter) JudgeMessage(ctx context.Context, msg *model.Message) bool {
	//  todo
	return true
}

// ReceiverOnline 判断接收方是否在线
func (mc *MessageCenter) ReceiverOnline(ctx context.Context, receiverId uint) bool {
	// todo change to cache
	_, exist := getNewsChannel(strconv.Itoa(int(receiverId)))
	return exist
}

func (mc *MessageCenter) SendMessageLater(ctx context.Context, msg *model.Message) error {
	// todo 触发通知
	go msgService.SaveMessage(ctx, msg)
	return nil
}
