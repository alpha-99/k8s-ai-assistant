package main

import (
	"fmt"
	"log"

	"github.com/api-agent/pkg/controllers"
	"github.com/api-agent/pkg/core/tools"
	"github.com/api-agent/pkg/models"
	"github.com/api-agent/pkg/services"
	"github.com/api-agent/pkg/sys"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func main() {
	sc := sys.InitConfig()

	var api models.OpenAPI
	err := yaml.Unmarshal([]byte(sc.Apis.API), &api)
	if err != nil {
		log.Fatalln(err)
	}

	tool, _ := tools.ParseOpenAPIToToolBundle(&api)
	fmt.Println("tools: ", tool)
	chatCompletionService := services.NewChatCompletionService(sc, tool)

	chatCompletionCtl := controllers.NewChatCompletionCtl(chatCompletionService)

	r := gin.New()

	r.POST("/v1/chat_messages", chatCompletionCtl.ChatCompletion())

	r.Run(":8080")

	//agent.Run(sc, tools, `帮我把"何以解忧，唯有暴富"翻译成英文`)
	//agent.Run(sc, tools, `济南奥体中心附近游泳馆有哪些`)
}
