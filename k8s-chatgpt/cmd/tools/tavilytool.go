package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TavilyTool
type TavilyTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

func NewTavilyTool() *TavilyTool {
	return &TavilyTool{
		Name:        "TavilyTool",
		Description: `Search the web for information on a topic`,
		ArgsSchema:  `description: "要搜索的内容，格式是字符串" example: "C罗是谁？"`,
	}
}

// RequestParams 定义请求参数
type RequestParams struct {
	APIKey                   string   `json:"api_key"`
	Query                    string   `json:"query"`
	SearchDepth              string   `json:"search_depth,omitempty"`
	Topic                    string   `json:"topic,omitempty"`
	Days                     int      `json:"days,omitempty"`
	MaxResults               int      `json:"max_results,omitempty"`
	IncludeImages            bool     `json:"include_images,omitempty"`
	IncludeImageDescriptions bool     `json:"include_image_descriptions,omitempty"`
	IncludeAnswer            bool     `json:"include_answer,omitempty"`
	IncludeRawContent        bool     `json:"include_raw_content,omitempty"`
	IncludeDomains           []string `json:"include_domains,omitempty"`
	ExcludeDomains           []string `json:"exclude_domains,omitempty"`
}

// Response 定义API响应的结构
type Response struct {
	Query        string         `json:"query"`
	Answer       string         `json:"answer,omitempty"`
	ResponseTime float64        `json:"response_time"`
	Images       []Image        `json:"images,omitempty"`
	Results      []SearchResult `json:"results,omitempty"`
}

// Image 定义图像结构
type Image struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// SearchResult 定义搜索结果结构
type SearchResult struct {
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Content       string  `json:"content"`
	RawContent    string  `json:"raw_content,omitempty"`
	Score         float64 `json:"score"`
	PublishedDate string  `json:"published_date,omitempty"`
}

type FinalResult struct {
	Title string
	Link  string
	//Snippet string
}

func (k *TavilyTool) Run(query string) ([]FinalResult, error) {
	url := "https://api.tavily.com/search"
	apiKey := "tvly-xxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	requestParam := RequestParams{
		APIKey:      apiKey,
		Query:       query,
		Days:        7,
		MaxResults:  7,
		SearchDepth: "basic",
	}

	// 初始化client
	client := http.Client{}
	body, err := json.Marshal(requestParam)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Json : %w", err)
	}

	// 创建新的HTTP请求
	req, err := http.NewRequest("post", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	finalResult := make([]FinalResult, 0)
	for _, res := range response.Results {
		finalResult = append(finalResult, FinalResult{
			Title: res.Title,
			Link:  " link: " + res.URL,
		})
	}

	return finalResult, nil
}
