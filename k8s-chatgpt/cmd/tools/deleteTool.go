package tools

import (
	"k8s-assistant/cmd/utils"
	"strings"
)

type DeleteToolParam struct {
	Resource  string `json:"resource"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// DeleteTool 表示一个工具，用于列出 k8s 资源命令。
type DeleteTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewDeleteTool 创建一个新的 DeleteTool 实例。
func NewDeleteTool() *DeleteTool {
	return &DeleteTool{
		Name:        "DeleteTool",
		Description: "用于删除指定命名空间的指定 Kubernetes 资源，例如删除某 pod 等等",
		ArgsSchema:  `{"type":"object","properties":{"resource":{"type":"string", "description": "指定的 k8s 资源类型，例如 pod, service等等"}, "name":{"type":"string", "description": "指定的某 k8s 资源实例的名称"}, "namespace":{"type":"string", "description": "指定的 k8s 资源所在命名空间"}}`,
	}
}

// Run 执行命令并返回输出。
func (d *DeleteTool) Run(resource, name, ns string) error {
	resource = strings.ToLower(resource)

	url := "http://localhost:8080/" + resource + "?ns=" + ns + "&name=" + name

	_, err := utils.DeleteHTTP(url)

	return err
}
