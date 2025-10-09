package tools

import (
	"strings"

	"github.com/k8s-chatgpt/cmd/utils"
)

type ListToolParam struct {
	Resource  string `json:"resource"`
	Namespace string `json:"namespace"`
}

// ListTool 表示一个工具，用于列出 k8s 资源命令。
type ListTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewListTool 创建一个新的 ListTool 实例。
func NewListTool() *ListTool {
	return &ListTool{
		Name:        "ListTool",
		Description: "用于列出指定命名空间的指定 Kubernetes 资源列表，例如 pod列表等等",
		ArgsSchema:  `{"type":"object","properties":{"resource":{"type":"string", "description": "指定的 k8s 资源类型，例如 pod, service等等"}, "namespace":{"type":"string", "description": "指定的 k8s 命名空间"}}`,
	}
}

// Run 执行命令并返回输出。
func (l *ListTool) Run(resource string, ns string) (string, error) {
	resource = strings.ToLower(resource)

	url := "http://localhost:8080/" + resource + "?ns=" + ns

	s, err := utils.GetHTTP(url)

	return s, err
}
