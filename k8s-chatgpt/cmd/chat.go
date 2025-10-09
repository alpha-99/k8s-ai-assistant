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

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		createTool := tools.NewCreateTool()
		listTool := tools.NewListTool()
		deleteTool := tools.NewDeleteTool()
		humanTool := tools.NewHumanTool()
		clusterTool := tools.NewClusterTool()

		scanner := bufio.NewScanner(cmd.InOrStdin())
		fmt.Println("你好，我是k8s助手，请问有什么可以帮助你？（请输入 'exit' 退出程序）：")
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

			prompt := buildPrompt(createTool, listTool, deleteTool, humanTool, clusterTool, input)
			ai.MessageStore.AddForUser(prompt)
			i := 1
			for {
				first_response := ai.NormalChat(ai.MessageStore.ToMessage())
				fmt.Printf("==========第%d轮回答==========\n", i)
				fmt.Println(first_response.Content)

				regexPattern := regexp.MustCompile(`Final Answer:\s*(.*)`)
				finalAnswer := regexPattern.FindStringSubmatch(first_response.Content)
				if len(finalAnswer) > 0 {
					fmt.Println("===========GPT 最终回复==============")
					fmt.Println(first_response.Content)
					break
				}

				ai.MessageStore.AddForAssistant(first_response.Content)

				regexAction := regexp.MustCompile(`Action:\s*(.*?)[\n]`)
				regexActionInput := regexp.MustCompile(`Action Input:\s*({[\s\S]*?})`)

				action := regexAction.FindStringSubmatch(first_response.Content)
				actionInput := regexActionInput.FindStringSubmatch(first_response.Content)
				if len(action) > 1 && len(actionInput) > 1 {
					i++
					Observation := "Observation: %s"
					switch action[1] {
					case createTool.Name:
						var param tools.CreateToolParam
						_ = json.Unmarshal([]byte(actionInput[1]), &param)

						output := createTool.Run(param.Prompt, param.Resource)
						Observation = fmt.Sprintf(Observation, output)
					case deleteTool.Name:
						var param tools.DeleteToolParam
						_ = json.Unmarshal([]byte(actionInput[1]), &param)

						err := deleteTool.Run(param.Resource, param.Name, param.Namespace)
						if err != nil {
							Observation = fmt.Sprintf(Observation, "删除失败")
						} else {
							Observation = fmt.Sprintf(Observation, "删除成功")
						}

					case listTool.Name:
						var param tools.ListToolParam
						_ = json.Unmarshal([]byte(actionInput[1]), &param)

						output, _ := listTool.Run(param.Resource, param.Namespace)
						Observation = fmt.Sprintf(Observation, output)
					case humanTool.Name:
						var param tools.HumanToolParam
						_ = json.Unmarshal([]byte(actionInput[1]), &param)

						output := humanTool.Run(param.Prompt)
						Observation = fmt.Sprintf(Observation, output)
					case clusterTool.Name:
						output, _ := clusterTool.Run()
						Observation = fmt.Sprintf(Observation, output)
					}

					prompt = first_response.Content + Observation
					fmt.Printf("========第%d轮的prompt========\n", i)
					fmt.Println(prompt)
					ai.MessageStore.AddForUser(prompt)
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func buildPrompt(createTool *tools.CreateTool, listTool *tools.ListTool, deleteTool *tools.DeleteTool, humanTool *tools.HumanTool, clustersTool *tools.ClusterTool, query string) string {
	createToolDef := "Name: " + createTool.Name + "\nDescription: " + createTool.Description + "\nArgsSchema: " + createTool.ArgsSchema + "\n"
	listToolDef := "Name: " + listTool.Name + "\nDescription: " + listTool.Description + "\nArgsSchema: " + listTool.ArgsSchema + "\n"
	deleteToolDef := "Name: " + deleteTool.Name + "\nDescription: " + deleteTool.Description + "\nArgsSchema: " + deleteTool.ArgsSchema + "\n"
	humanToolDef := "Name: " + humanTool.Name + "\nDescription: " + humanTool.Description + "\nArgsSchema: " + humanTool.ArgsSchema + "\n"
	clusterToolDef := "Name: " + clustersTool.Name + "\nDescription: " + clustersTool.Description + "\n"

	toolsList := make([]string, 0)
	toolsList = append(toolsList, createToolDef, listToolDef, deleteToolDef, humanToolDef, clusterToolDef)

	tool_names := make([]string, 0)
	tool_names = append(tool_names, createTool.Name, listTool.Name, deleteTool.Name, humanTool.Name, clustersTool.Name)

	prompt := fmt.Sprintf(promptTpl.Template, toolsList, tool_names, "", query)

	return prompt
}
