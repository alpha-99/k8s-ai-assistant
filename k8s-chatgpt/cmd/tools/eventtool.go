package tools

import (
	"fmt"

	"github.com/k8s-chatgpt/cmd/utils"
)

type EventToolParam struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
}

type EventTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

func NewEventTool() *EventTool {
	return &EventTool{
		Name:        "EventTool",
		Description: "用于查看 k8s pod 的 event 事件",
		ArgsSchema:  `{"type":"object","properties":{"podName":{"type":"string", "description": "指定的 pod 名称"}, "namespace":{"type":"string", "description": "指定的 k8s 命名空间"}}`,
	}
}

// Run 执行命令并返回输出。
func (e *EventTool) Run(podName string, ns string) (string, error) {
	url := "http://localhost:8080/pods/events" + "?ns=" + ns + "&podname=" + podName
	fmt.Println("url: ", url)
	s, err := utils.GetHTTP(url)

	return s, err
}
