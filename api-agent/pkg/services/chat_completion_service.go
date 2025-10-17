package services

import (
	"github.com/api-agent/pkg/core/agent"
	"github.com/api-agent/pkg/models"
)

type ChatCompletionService struct {
	sc    *models.Config
	tools []models.ApiToolBundle
}

func NewChatCompletionService(sc *models.Config, tools []models.ApiToolBundle) *ChatCompletionService {
	return &ChatCompletionService{
		sc:    sc,
		tools: tools,
	}
}

func (s *ChatCompletionService) ChatCompletion(query string) (string, error) {
	return agent.Run(s.sc, s.tools, query)
}
