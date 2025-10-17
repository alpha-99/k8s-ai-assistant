package outputparser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ParseActionFromText 解析输入字符串中的 action 和 action_input
func HandleReActOutput(input string) (string, map[string]interface{}, error) {
	// 定义匹配 ``` ``` 中内容的正则表达式
	// 使用非贪婪模式 .*? 以匹配最短的内容
	re := regexp.MustCompile("(?s)```(.*?)```")
	//re := regexp.MustCompile("(?s)```(\\w*)\\s*(.*?)```")

	// 查找所有匹配的代码块
	matches := re.FindStringSubmatch(input)
	if len(matches) == 0 {
		return "", nil, fmt.Errorf("未找到被 ``` ``` 包围的 JSON 内容")
	}

	// 假设第一个匹配的代码块是我们需要的 JSON
	jsonString := strings.TrimSpace(matches[1])
	fmt.Println("提取的 JSON 字符串:")
	fmt.Println(jsonString)

	// 反序列化 JSON
	var actionData map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &actionData)
	if err != nil {
		fmt.Println("JSON 反序列化失败:", err)
		return "", nil, err
	}

	fmt.Println("反序列化后的 JSON 数据: ", actionData)

	// 提取 action 和 action_input
	action, ok := actionData["action"].(string)
	if !ok {
		fmt.Println("JSON 中缺少 'action' 字段，或类型不匹配")
		return "", nil, err
	}

	// 检查 action 是否是 "Final Answer"
	if action == "Final Answer" {
		actionInput, ok := actionData["action_input"].(string)
		if !ok {
			fmt.Println("JSON 中缺少 'action_input' 字段，或类型不匹配")
			return "", nil, err
		}

		actionInputMap := map[string]interface{}{
			"finalAnswer": actionInput,
		}

		return action, actionInputMap, nil
	}

	actionInput, ok := actionData["action_input"].(map[string]interface{})

	if !ok {
		fmt.Println("JSON 中缺少 'action_input' 字段，或类型不匹配")
		return "", nil, err
	}

	return action, actionInput, nil

}
