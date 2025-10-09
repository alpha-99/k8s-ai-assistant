/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/k8s-chatgpt/cmd/ai"
	"github.com/k8s-chatgpt/cmd/promptTpl"
	"github.com/k8s-chatgpt/cmd/tools"

	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logTool := tools.NewLogTool()
		logToolDef := "Name: " + logTool.Name + "\nDescription: " + logTool.Description + "\nArgsSchema: " + logTool.ArgsSchema + "\n"

		eventTool := tools.NewEventTool()
		eventToolDef := "Name: " + eventTool.Name + "\nDescription: " + eventTool.Description + "\nArgsSchema: " + eventTool.ArgsSchema + "\n"

		toolsList := make([]string, 0)
		toolsList = append(toolsList, logToolDef, eventToolDef)

		tool_names := make([]string, 0)
		tool_names = append(tool_names, logTool.Name, eventTool.Name)

		scanner := bufio.NewScanner(cmd.InOrStdin())
		fmt.Println("你好，我是k8s助手，请问有什么可以帮你？（输入 'exit' 退出程序）:")
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			input := scanner.Text()
			if input == "exit" {
				fmt.Println("再见！")
				return
			}

			prompt := fmt.Sprintf(promptTpl.Template, toolsList, tool_names, "", input)

			// 注入用户的prompt
			ai.MessageStore.AddForUser(prompt)
			i := 1
			for {
				firstResponse := ai.NormalChat(ai.MessageStore.ToMessage())
				fmt.Printf("==============第%d轮回答=============\n", i)
				fmt.Println(firstResponse.Content)

				regexPattern := regexp.MustCompile(`Final Answer:\s*(.*)`)
				finalAnswer := regexPattern.FindStringSubmatch(firstResponse.Content)
				if len(finalAnswer) > 1 {
					fmt.Println("========最终 GPT 回复========")
					fmt.Println(firstResponse.Content)
					break
				}

				ai.MessageStore.AddForAssistant(firstResponse.Content)

				regexAction := regexp.MustCompile(`Action:\s*(.*?)[\n]`)
				regexActionInput := regexp.MustCompile(`Action Input:\s*({[\s\S]*?})`)

				action := regexAction.FindStringSubmatch(firstResponse.Content)
				actionInput := regexActionInput.FindStringSubmatch(firstResponse.Content)

				if len(action) > 0 && len(actionInput) > 0 {
					i++
					observation := "Observation: %s"
					switch action[1] {
					case logTool.Name:
						var param tools.LogToolParam
						_ = json.Unmarshal([]byte(actionInput[1]), &param)

						output, _ := logTool.Run(param.PodName, param.Namespace)
						observation = fmt.Sprintf(observation, output)
					case eventTool.Name:
						var param tools.EventToolParam
						_ = json.Unmarshal([]byte(actionInput[1]), &param)

						output, _ := eventTool.Run(param.PodName, param.Namespace)
						observation = fmt.Sprintf(observation, output)
					}

					prompt = firstResponse.Content + observation
					fmt.Printf("========第%d轮的prompt========\n", i)
					fmt.Println(prompt)
					ai.MessageStore.AddForUser(prompt)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyzeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analyzeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
