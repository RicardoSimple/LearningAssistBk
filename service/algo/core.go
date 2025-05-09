package algo

import (
	"encoding/json"
	"fmt"
	"learning-assistant/util"
	"learning-assistant/util/log"
)

type Request struct {
	Model            string         `json:"model"`
	Messages         []ChatMessage  `json:"messages"`
	FrequencyPenalty float64        `json:"frequency_penalty"`
	PresencePenalty  float64        `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
	MaxTokens        uint           `json:"max_tokens"`
	Stop             []string       `json:"stop"`
	Stream           bool           `json:"stream"`
	Temperature      float64        `json:"temperature"`
}
type Response struct {
	Id                string    `json:"id"`
	Object            string    `json:"object"`
	Created           uint64    `json:"created"`
	Model             string    `json:"model"`
	SystemFingerprint string    `json:"system_fingerprint"`
	Choices           []Choice  `json:"choices"`
	Usage             ChatUsage `json:"usage"`
}

func (rq Request) DefaultChatParam() {
	rq.Model = ChatModel
	rq.FrequencyPenalty = 0.0
	rq.MaxTokens = 4096
	rq.PresencePenalty = 0.0
	rq.ResponseFormat.Type = ResponseTypeText
	rq.Temperature = 1
}

func (c *ChatClient) Chat(messages []ChatMessage, model string) (string, error) {
	url := c.BaseUrl + "/v1/chat/completions"
	req := Request{
		Model:    model,
		Messages: messages,
	}
	req.DefaultChatParam()
	var resp Response
	err := util.DoJsonPost(url, map[string]string{
		"Authorization": "Bearer " + c.ApiKey,
	}, req, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("响应为空")
	}
	return resp.Choices[0].Delta.Content, nil
}

func (c *ChatClient) ChatStream(
	messages []ChatMessage,
	model string,
	onMessage func(content string),
) error {
	url := c.BaseUrl + "/v1/chat/completions"
	req := Request{
		Model:    model,
		Messages: messages,
		Stream:   true, // 添加 Stream 参数
	}
	req.DefaultChatParam()
	log.Info("send chat stream:" + messages[len(messages)-1].Content)
	return util.DoJsonPostStream(url, map[string]string{
		"Authorization": "Bearer " + c.ApiKey,
	}, req, func(data string) {
		// 解析每次推送的 JSON 包
		log.Info(data)
		var partial Response
		if err := json.Unmarshal([]byte(data), &partial); err == nil {
			if len(partial.Choices) > 0 {
				log.Info("chat response: " + partial.Choices[0].Delta.Content)
				onMessage(partial.Choices[0].Delta.Content)
			}
		}
	})
}
