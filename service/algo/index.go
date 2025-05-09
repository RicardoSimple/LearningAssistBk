package algo

import (
	"bufio"
	"fmt"
	"learning-assistant/util/log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ChatModel     = "deepseek-chat"
	ReasonerModel = "deepseek-reasoner"

	SystemRole    = "system"
	UserRole      = "user"
	AssistantRole = "assistant"

	ResponseTypeText = "text"
	ResponseTypeJson = "json_object"
)

type ChatClient struct {
	BaseUrl string `json:"base_url"`
	ApiKey  string `json:"api_key"`
}

var client *ChatClient

func createClient() *ChatClient {
	// 读取baseurl
	baseUrlPath := filepath.Join(".", "chat_base_url")
	baseUrlFile, err := os.Open(baseUrlPath)
	if err != nil {
		panic(fmt.Sprintf("读取baseurl文件失败: %v", err))
	}
	defer baseUrlFile.Close()
	scanner := bufio.NewScanner(baseUrlFile)
	lines := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
			log.Debug(line)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("读取api配置文件失败: %v", err))
	}

	if len(lines) < 2 {
		panic("配置文件格式错误，需至少包含 base_url 和 api_key 两行")
	}
	return &ChatClient{
		BaseUrl: lines[0],
		ApiKey:  lines[1],
	}
}

func GetClient() *ChatClient {
	if client == nil {
		client = createClient()
	}
	return client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ResponseFormat struct {
	Type string `json:"type"`
}

type ChatUsage struct {
	PromptTokens          uint                `json:"prompt_tokens"`
	CompletionTokens      uint                `json:"completion_tokens"`
	TotalTokens           uint                `json:"total_tokens"`
	PromptTokensDetails   PromptTokensDetails `json:"prompt_tokens_details"`
	PromptCacheHitTokens  uint                `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens uint                `json:"prompt_cache_miss_tokens"`
}
type PromptTokensDetails struct {
	CachedTokens uint `json:"cached_tokens"`
}

type Choice struct {
	Index        uint        `json:"index"`
	Message      ChatMessage `json:"message"`
	Delta        ChoiceDelta `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}
type ChoiceDelta struct {
	Content string `json:"content"`
}
