package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// KubeInput表示KubeTool的输入
type KubeInput struct {
	Commands string
}

// KubeTool表示一个工具，用于运行kubernetes命令
type KubeTool struct {
	Name        string
	Description string
	ArgsSchema  KubeInput
}

// NewKubeTool创建一个新的KubeTool实例
func NewKubeTool() *KubeTool {
	return &KubeTool{
		Name:        "KubeTool",
		Description: "用于在 Kubernetes 集群上运行 k8s 相关命令（kubectl、helm）的工具。",
		ArgsSchema:  KubeInput{Commands: `description: "要运行的 kubectl/helm 相关命令。" example: "kubectl get pods"`},
	}
}

// Run执行命令并返回输出
func (k *KubeTool) Run(commands string) (string, error) {
	parsedCommands := k.parseCommands(commands)

	splitedCommands := k.splitCommands(parsedCommands)

	// 在这里，您通常会使用 os/exec 包来执行命令并返回输出。
	cmd := exec.Command(splitedCommands[0], splitedCommands[1:]...)

	// 运行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}

	return fmt.Sprintf("运行结果：%s", output), nil
}

// parseCommands 清理命令字符串
func (k *KubeTool) parseCommands(commands string) string {
	return strings.TrimSpace(strings.Trim(commands, "\"`"))
}

// splitCommands切割输入的命令字符串
func (k *KubeTool) splitCommands(commands string) []string {
	return strings.Split(commands, " ")
}
