package tools

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RequestGet结构体，用于处理HTTP请求
type RequestTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

func NewRequestTool() *RequestTool {
	return &RequestTool{
		Name: "RequestTool",
		Description: `
		A portal to the internet. Use this when you need to get specific
    content from a website. Input should be a url (i.e. https://www.kubernetes.io/releases).
    The output will be the text response of the GET request.
		`,
		ArgsSchema: `description: "要访问的website，格式是字符串" example: "https://www.kubernetes.io/releases"`,
	}
}

func (r *RequestTool) Run(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取URL失败，%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return r.parseHTML(string(body)), nil
}

func (r *RequestTool) parseHTML(body string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return ""
	}

	// 移除不需要的标签
	doc.Find("header, footer, script, style").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	return doc.Find("body").Text()
}
