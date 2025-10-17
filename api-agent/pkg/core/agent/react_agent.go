package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	outputparser "github.com/api-agent/pkg/core/agent/output_parser"
	prompttemplate "github.com/api-agent/pkg/core/agent/template"
	"github.com/api-agent/pkg/core/ai"
	"github.com/api-agent/pkg/core/tools"
	"github.com/api-agent/pkg/models"
)

func organizeToolsPrompt(tools []models.ApiToolBundle) (string, []string) {
	messageTools := make([]models.PromptMessageTool, 0)
	toolNames := []string{}

	for _, tool := range tools {
		toolNames = append(toolNames, tool.OperationId)

		properties := make([]models.Properties, 0)
		for _, parameter := range tool.Parameters {
			properties = append(properties, models.Properties{
				Type:        parameter.Type,
				Name:        parameter.Name,
				Description: parameter.LLMDescription,
				Enum:        parameter.Enum,
			})
		}

		required := make([]string, 0)
		for _, parameter := range tool.Parameters {
			if parameter.Required {
				required = append(required, parameter.Name)
			}
		}

		parameters := models.Parameters{
			Type:       "object",
			Properties: properties,
			Required:   required,
		}

		messageTool := models.PromptMessageTool{
			Name:        tool.OperationId,
			Description: tool.Summary,
			Parameters:  parameters,
		}

		messageTools = append(messageTools, messageTool)
	}

	jsonMessageTools, err := json.Marshal(messageTools)
	if err != nil {
		panic(err)
	}

	return string(jsonMessageTools), toolNames
}

func organizeHistoricMessage() string {
	return ""
}

func organizeReActTemplate(instruction string, tools []models.ApiToolBundle, query string) string {
	messageTools, messageToolNames := organizeToolsPrompt(tools)

	historicMessage := organizeHistoricMessage()

	// 填充数据
	data := models.TemplateData{
		Instruction:      instruction,
		Tools:            messageTools,
		ToolNames:        messageToolNames,
		HistoricMessages: historicMessage,
		Query:            query,
	}

	// 加载渲染模板
	tmpl, err := template.New("prompt").Parse(prompttemplate.EN_Template)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// 创建一个缓冲区来接收模板渲染输出
	var result bytes.Buffer

	// 执行模板并填充数据
	err = tmpl.Execute(&result, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return result.String()
}

func Run(sc *models.Config, toolBundles []models.ApiToolBundle, query string) (string, error) {
	prompt := organizeReActTemplate(sc.Instruction, toolBundles, query)
	ai.MessageStore.AddForUser(prompt)

	var action string
	var actionInput map[string]interface{}
	var err error

	var iterationSteps = 1
	for {
		firstResponse := ai.NormalChat(ai.MessageStore.ToMessage())
		fmt.Printf("========第%d轮回答========\n", iterationSteps)
		fmt.Println(firstResponse.Content)

		ai.MessageStore.AddForAssistant(firstResponse.Content)

		action, actionInput, err = outputparser.HandleReActOutput(firstResponse.Content)
		if err != nil {
			fmt.Println("Error: ", err)
			return "", err
		}

		fmt.Println("Action:", action)
		fmt.Println("Action_Input:", actionInput)

		if action == "Final Answer" {
			break
		}

		if iterationSteps >= sc.MaxIterationSteps {
			actionInput["finalAnswer"] = "已超出允许的最大迭代次数"
			break
		}

		iterationSteps++
		observation := "Observation: %s"
		for _, toolBundle := range toolBundles {
			if action == toolBundle.OperationId {
				respBody, statusCode, err := tools.ToolInvoke(sc.Apis.APIProvider.APIKey, toolBundle.Method, toolBundle.ServerURL, toolBundle, actionInput)
				if err != nil {
					fmt.Println("Error: ", err)
					return "", err
				}

				fmt.Println("StatusCode:", statusCode)
				fmt.Println("Response:", string(respBody))
				observation = fmt.Sprintf(observation, string(respBody))
				break
			}
		}
		ai.MessageStore.AddForUser(observation)
	}

	return actionInput["finalAnswer"].(string), nil
}
