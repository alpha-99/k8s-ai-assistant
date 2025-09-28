package tools

import (
	"encoding/json"
	"fmt"
	"k8s-assistant/cmd/ai"
	"k8s-assistant/cmd/promptTpl"
	"k8s-assistant/cmd/utils"

	"github.com/sashabaranov/go-openai"
)

type CreateTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

func NewCreateTool() *CreateTool {
	return &CreateTool{
		Name:        "CreateTool",
		Description: "用于在指定命名空间创建指定 Kubernetes 资源，例如创建某 pod 等等",
		ArgsSchema:  `{"type":"object","properties":{"prompt":{"type":"string", "description": "把用户提出的创建资源的prompt原样放在这，不要做任何改变"},"resource":{"type":"string", "description": "指定的 k8s 资源类型，例如 pod, service等等"}}}`,
	}
}

// 定义结构体解析json响应
type response struct {
	Data string `json:"data"`
}

type CreateToolParam struct {
	Prompt   string `json:"prompt"`
	Resource string `json:"resource"`
}

// Run 执行命令并返回
func (c *CreateTool) Run(prompt string, resource string) string {
	// 让大模型生成yaml
	messages := make([]openai.ChatCompletionMessage, 2)

	messages[0] = openai.ChatCompletionMessage{Role: "system", Content: promptTpl.SystemPrompt}
	messages[1] = openai.ChatCompletionMessage{Role: "user", Content: prompt}

	resp := ai.NormalChat(messages)
	fmt.Println("-----------------------")
	fmt.Println(resp.Content)

	// 创建json对象 {"yaml": "xxx"}
	body := map[string]string{"yaml": resp.Content}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}

	url := "http://localhost:8080/" + resource
	s, err := utils.PostHTTP(url, jsonBody)
	if err != nil {
		return err.Error()
	}

	var response response
	err = json.Unmarshal([]byte(s), &response)
	if err != nil {
		return err.Error()
	}

	return response.Data
}
