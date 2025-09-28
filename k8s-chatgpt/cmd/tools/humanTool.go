package tools

import "fmt"

type HumanToolParam struct {
	Prompt string `json:"prompt"`
}

// HumanTool 表示一个工具，用于列出 k8s 资源命令。
type HumanTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewDeleteTool 创建一个新的 HumanTool 实例。
func NewHumanTool() *HumanTool {
	return &HumanTool{
		Name:        "HumanTool",
		Description: "当你判断出要执行一些不可逆的危险操作时，比如删除动作，需要先用本工具向人类发起确认",
		ArgsSchema:  `{"type":"object","properties":{"prompt":{"type":"string", "description": "你要向人类寻求帮助的内容", "example": "请确认是否要删除 default 命名空间下的 foo-app pod"}}}`,
	}
}

// Run 执行命令并返回输出。
func (d *HumanTool) Run(prompt string) string {
	fmt.Print(prompt, " ")
	var input string
	fmt.Scanln(&input)
	return input
}
