package template

const EN_Template = `
Respond to the human as helpfully and accurately as possible.

{{.Instruction}}

You have access to the following tools:

{{.Tools}}

Use a json blob to specify a tool by providing an action key (tool name) and an action_input key (tool input).
Return the json blob directly after the backticks, without adding the word 'json' after the backticks.
Valid "action" values: "Final Answer" or {{.ToolNames}}

Provide only ONE action per $JSON_BLOB, as shown:
` + "```" + `
{
  "action": $TOOL_NAME,
  "action_input": $ACTION_INPUT
}
` + "```" + `
Follow this format:
Question: input question to answer
Thought: consider previous and subsequent steps 
Action: ` + "```" + `$JSON_BLOB` + "```" + `

Observation: action result 
... (repeat Thought/Action/Observation N times)
Thought: I know what to respond
Action:` + "```" + `
{
  "action": "Final Answer",
  "action_input": "Final response to human"
}
` + "```" + `
Begin! Reminder to ALWAYS respond with a valid json blob of a single action. Use tools if necessary. Respond directly if appropriate.Format is Action:` + "```" + `$JSON_BLOB` + "```" + `then Observation:.

{{.HistoricMessages}}

Question: {{.Query}}
`
