package tools

import "github.com/k8s-chatgpt/cmd/utils"

type LogToolParam struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
}

type LogTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

func NewLogTool() *LogTool {
	return &LogTool{
		Name:        "LogTool",
		Description: "用于查看 k8s pod 的 log 日志",
		ArgsSchema:  `{"type":"object","properties":{"podName":{"type":"string", "description": "指定的 pod 名称"}, "namespace":{"type":"string", "description": "指定的 k8s 命名空间"}}`,
	}
}

func (l *LogTool) Run(podName string, ns string) (string, error) {
	url := "http://localhost:8080/pods/logs" + "?ns=" + ns + "&podname=" + podName

	s, err := utils.GetHTTP(url)

	return s, err
}
