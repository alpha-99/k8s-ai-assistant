package controllers

import (
	"github.com/api-agent/pkg/models"
	"github.com/api-agent/pkg/services"
	"github.com/gin-gonic/gin"
)

type ChatCompletionCtl struct {
	chatCompletionService *services.ChatCompletionService
}

func NewChatCompletionCtl(svc *services.ChatCompletionService) *ChatCompletionCtl {
	return &ChatCompletionCtl{chatCompletionService: svc}
}

func (chat *ChatCompletionCtl) ChatCompletion() func(c *gin.Context) {
	return func(c *gin.Context) {
		var message models.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(400, gin.H{"error": "解析请求体失败: " + err.Error()})
		}

		response, err := chat.chatCompletionService.ChatCompletion(message.Message)
		if err != nil {
			c.JSON(400, gin.H{"error": "询问失败: " + err.Error()})
		}

		c.JSON(200, gin.H{"message": response})
	}
}
