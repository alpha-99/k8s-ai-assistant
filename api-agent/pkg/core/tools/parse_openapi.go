package tools

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/api-agent/pkg/models"
	"github.com/google/uuid"
)

// 定义错误类型
var ErrNoServerFound = errors.New("no server found in the swagger yaml")
var ErrNoPathFound = errors.New("no path found in the swagger yaml")
var ErrNoOperationId = errors.New("no operationId found in operation")
var ErrNoSummaryOrDescription = errors.New("no summary or description found in operation")

type ToolApiSchemaError struct {
	Message string
}

func (e *ToolApiSchemaError) Error() string {
	return e.Message
}

func ParseSwaggerToOpenAPI(swagger *models.Swagger) (*models.OpenAPI, error) {
	// 从swagger获取信息
	info := swagger.Info
	if info == nil {
		info = map[string]string{"title": "Swagger", "description": "Swagger", "version": "1.0.0"}
	}

	servers := swagger.Servers
	if len(servers) == 0 {
		return nil, &ToolApiSchemaError{Message: ErrNoServerFound.Error()}
	}

	// 初始化 OpenAPI对象
	openAPI := &models.OpenAPI{
		OpenAPI: "3.0.0",
		Info: map[string]string{
			"title":       info["title"],
			"description": info["description"],
			"version":     info["version"],
		},
		Servers:    swagger.Servers,
		Paths:      make(map[string]map[string]interface{}),
		Components: map[string]map[string]interface{}{"schemas": make(map[string]interface{})},
	}

	// 检查paths是否存在
	if len(swagger.Paths) == 0 {
		return nil, &ToolApiSchemaError{Message: ErrNoPathFound.Error()}
	}

	// 转换paths，swagger的path格式可以参考：https://swagger.org.cn/docs/specification/v2_0/paths-and-operations/
	for path, pathItem := range swagger.Paths {
		openAPI.Paths[path] = make(map[string]interface{})
		for method, operation := range pathItem {
			// 将operation断言为map[string]interface{}
			ops, ok := operation.(map[string]interface{})
			if !ok {
				return nil, &ToolApiSchemaError{Message: fmt.Sprintf("Invalid operation format for %s %s", method, path)}
			}

			callMethod, ok := ops["operationId"].(string)
			if !ok {
				return nil, &ToolApiSchemaError{Message: fmt.Sprintf("%s %s %s", ErrNoOperationId.Error(), method, path)}
			}

			// 如果没有summary或desciption，添加告警
			summary, hasSummary := ops["summary"].(string)
			description, hasDescription := ops["description"].(string)
			if !hasSummary || !hasDescription {
				log.Warnf("No summary or description found in operation %s %s.", method, path)
			}

			openAPI.Paths[path][method] = map[string]interface{}{
				"operationId": callMethod,
				"summary":     summary,
				"description": description,
				"parameters":  ops["parameters"],
				"responses":   ops["responses"],
			}

			// 检查是否存在requestBody
			// 文档格式参考：https://swagger.io/docs/specification/v3_0/describing-request-body/describing-request-body/?sbsearch=requestBody
			if requestBody, ok := ops["requestBody"]; ok {
				openAPI.Paths[path][method].(map[string]interface{})["requestBody"] = requestBody
			}
		}
	}

	// 转换definitions
	for name, definition := range swagger.Definitions {
		openAPI.Components["schemas"][name] = definition
	}

	return openAPI, nil
}

// 从OpenAPI转为ToolBundle
func ParseOpenAPIToToolBundle(openAPI *models.OpenAPI) ([]models.ApiToolBundle, error) {
	serverURL := openAPI.Servers[0]["url"].(string)

	// 列出所有的接口
	var interfaces []map[string]interface{}
	for path, pathItems := range openAPI.Paths {
		methods := []string{"get", "post", "delete", "patch", "head", "options", "trace"}
		for _, method := range methods {
			if methodItem, ok := pathItems[method]; ok {
				interfaces = append(interfaces, map[string]interface{}{
					"path":      path,
					"method":    method,
					"operation": methodItem,
				})
			}
		}
	}

	// 获取所有的参数并构建工具bundle
	var bundles []models.ApiToolBundle
	for _, iface := range interfaces {
		parameters := []models.ToolParameter{}
		operation := iface["operation"].(map[string]interface{})

		// 处理参数
		if params, ok := operation["parameters"].([]interface{}); ok {
			for _, param := range params {
				paramMap := param.(map[string]interface{})
				toolParam := models.ToolParameter{
					Name:           paramMap["name"].(string),
					Type:           "string", // 默认类型
					Required:       paramMap["required"].(bool),
					LLMDescription: paramMap["description"].(string),
					Default:        getDefault(paramMap),
				}

				// 类型处理
				if t := getParameterType(paramMap); t != "" {
					toolParam.Type = t
				}

				parameters = append(parameters, toolParam)
			}
		}

		// 处理请求体
		if requestBody, ok := operation["requestBody"].(map[string]interface{}); ok {
			if content, ok := requestBody["content"].(map[string]interface{}); ok {
				for _, contentType := range content {
					if bodySchema, ok := contentType.(map[string]interface{})["schema"].(map[string]interface{}); ok {
						required := bodySchema["required"].([]interface{})
						properties := bodySchema["properties"].(map[string]interface{})
						for name, prop := range properties {
							propMap := prop.(map[string]interface{})
							// 处理引用，如果有
							if ref, ok := propMap["$ref"].(string); ok {
								root := openAPI.Components["schemas"]
								segments := strings.Split(ref, "/")[1:]

								lastSegment := segments[len(segments)-1]
								propMap = root[lastSegment].(map[string]interface{})
							}

							toolParam := models.ToolParameter{
								Name:           name,
								Type:           "string", // 默认类型
								Required:       contains(required, name),
								LLMDescription: propMap["description"].(string),
								Default:        propMap["default"],
							}

							// 如果参数包含enum，则添加枚举值
							if enum, ok := propMap["enum"].([]interface{}); ok {
								var enumValues []string
								for _, e := range enum {
									enumValues = append(enumValues, e.(string))
								}
								toolParam.Enum = enumValues
							}

							// 类型处理
							if t := getParameterType(propMap); t != "" {
								toolParam.Type = t
							}

							parameters = append(parameters, toolParam)
						}
					}
				}
			}
		}

		// 检查参数是否重复
		paramCount := make(map[string]int)
		for _, param := range parameters {
			paramCount[param.Name]++
		}

		for name, count := range paramCount {
			if count > 1 {
				log.Warnf("Parameter %s is duplicated.", name)
			}
		}

		// 设置operationId
		if _, ok := operation["operationId"]; ok {
			// 如果没有operationId，使用path和method生成
			path := iface["path"].(string)
			path = strings.TrimPrefix(path, "/")
			// if strings.HasPrefix(path, "/") {
			// 	path = path[1:]
			// }

			// 移除特殊字符以确保operationId合法
			re := regexp.MustCompile("[^a-zA-Z0-9_-]")
			path = re.ReplaceAllString(path, "")
			if path == "" {
				path = uuid.New().String()
			}
			operation["operationId"] = fmt.Sprintf("%s_%s", path, iface["method"].(string))
		}

		// 构建ApiToolBundle
		bundles = append(bundles, models.ApiToolBundle{
			ServerURL:   serverURL + iface["path"].(string),
			Method:      iface["method"].(string),
			Summary:     getStringOrDefault(operation["description"], ""),
			OperationId: operation["operationId"].(string),
			Parameters:  parameters,
			OpenAPI:     operation,
		})
	}

	return bundles, nil
}

func getStringOrDefault(value interface{}, defaultValue string) string {
	if value == nil {
		return defaultValue
	}

	strValue, ok := value.(string)
	if !ok {
		return defaultValue
	}
	return strValue
}

func getDefault(param map[string]interface{}) interface{} {
	if schema, ok := param["schema"].(map[string]interface{}); ok {
		return schema["default"]
	}

	return nil
}

func getParameterType(param map[string]interface{}) string {
	// 根据实际情况来获取返回的正确类型
	if schema, ok := param["schema"].(map[string]interface{}); ok {
		if t, ok := schema["type"].(string); ok {
			return t
		}
	}

	return ""
}

// 判断元素是否在数组中
func contains(slice []interface{}, ele string) bool {
	for _, item := range slice {
		if item == ele {
			return true
		}
	}
	return false
}
