package tools

import "k8s-assistant/cmd/utils"

type ClusterTool struct {
	Name        string
	Description string
}

func NewClusterTool() *ClusterTool {
	return &ClusterTool{
		Name:        "ClusterTool",
		Description: "用于列出集群列表",
	}
}

// Run 执行命令并返回输出。
func (l *ClusterTool) Run() (string, error) {

	url := "http://localhost:8081/clusters"

	s, err := utils.GetHTTP(url)

	return s, err
}
