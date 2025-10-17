package models

// OpenAPI 和 Swagger的结构体定义
type Swagger struct {
	Info        map[string]string                 `json:"info"`
	Servers     []map[string]interface{}          `json:"servers"`
	Paths       map[string]map[string]interface{} `json:"paths"`
	Definitions map[string]interface{}            `json:"definitions"`
}

type OpenAPI struct {
	OpenAPI    string                            `json:"openapi"`
	Info       map[string]string                 `json:"info"`
	Servers    []map[string]interface{}          `json:"servers"`
	Paths      map[string]map[string]interface{} `json:"paths"`
	Components map[string]map[string]interface{} `json:"components"`
}

type ToolParameter struct {
	Name           string
	Type           string
	Required       bool
	LLMDescription string
	Default        interface{}
	Enum           []string
}

type ApiToolBundle struct {
	ServerURL   string
	Method      string
	Summary     string
	OperationId string
	Parameters  []ToolParameter
	OpenAPI     map[string]interface{}
}

type Properties struct {
	Type        string
	Name        string
	Description string
	Enum        []string
}

type Parameters struct {
	Type       string
	Properties []Properties
	Required   []string
}

type PromptMessageTool struct {
	Name        string
	Description string
	Parameters  Parameters
}

type TemplateData struct {
	Instruction      string
	Tools            string
	ToolNames        []string
	HistoricMessages string
	Query            string
}

type ReActOutput struct {
	Action      string `json:"action"`
	ActionInput string `json:"action_input"`
}

type ChatMessage struct {
	Message string `json:"message"`
}
