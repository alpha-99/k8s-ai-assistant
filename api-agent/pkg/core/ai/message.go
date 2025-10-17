package ai

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

var MessageStore ChatMessages

type ChatMessages []ChatMessage

type ChatMessage struct {
	Msg openai.ChatCompletionMessage
}

func init() {
	MessageStore = make(ChatMessages, 0)
	MessageStore.Clear() // 清理和初始化
}

func (cm *ChatMessages) Clear() {
	*cm = make(ChatMessages, 0)
	cm.AddForSystem("你是一个数学家")
}

func (cm *ChatMessages) AddForSystem(msg string) {
	cm.AddFor(msg, RoleSystem)
}

func (cm *ChatMessages) AddFor(msg, role string) {
	*cm = append(*cm, ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:    role,
			Content: msg,
		},
	})
}

func (cm *ChatMessages) AddForAssistant(msg string) {
	cm.AddFor(msg, RoleAssistant)
}

func (cm *ChatMessages) AddForUser(msg string) {
	cm.AddFor(msg, RoleUser)
}

// 组装prompt
func (cm *ChatMessages) ToMessage() []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, len(*cm))

	for index, c := range *cm {
		messages[index] = c.Msg
	}

	return messages
}

func (cm *ChatMessages) GetLast() string {
	if len(*cm) == 0 {
		return "Got nothing"
	}

	return (*cm)[len(*cm)-1].Msg.Content
}

func NewOpenAiClient() *openai.Client {
	token := os.Getenv("ALI_API_KEY")
	dashscope_url := "https://dashscope.aliyuncs.com/compatible-mode/v1"

	config := openai.DefaultConfig(token)
	config.BaseURL = dashscope_url

	return openai.NewClientWithConfig(config)
}

// chat对话
func NormalChat(message []openai.ChatCompletionMessage) openai.ChatCompletionMessage {
	c := NewOpenAiClient()
	resp, err := c.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model:    "qwen-max",
		Messages: message,
	})

	if err != nil {
		log.Println(err)
		return openai.ChatCompletionMessage{}
	}

	return resp.Choices[0].Message
}
