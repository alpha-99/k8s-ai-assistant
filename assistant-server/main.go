/*
Copyright © 2025 k8s-assistant

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"k8s-assistant/pkg/config"
	"k8s-assistant/pkg/controllers"
	"k8s-assistant/pkg/service"

	"github.com/gin-gonic/gin"
)

func main() {
	runServer()
}

func runServer() {
	k8sconfig := config.NewK8sConfig().InitRestConfig()
	restMapper := k8sconfig.InitRestMapper()
	dynamicClient := k8sconfig.InitDynamicClient()
	informer := k8sconfig.InitInformer()

	clientSet := k8sconfig.InitClientSet()

	resourceCtl := controllers.NewResourceCtl(service.NewResourceService(&restMapper, dynamicClient, informer))
	podLogCtl := controllers.NewPodLogEventCtl(service.NewPodLogEventService(clientSet))

	r := gin.New()
	r.GET("/:resource", resourceCtl.List())
	r.DELETE("/:resource", resourceCtl.Delete())
	r.POST("/:resource", resourceCtl.Create())
	r.GET("/get/gvr", resourceCtl.GetGVR())

	r.GET("/pods/logs", podLogCtl.GetLog())
	r.GET("/pods/events", podLogCtl.GetEvent())

	r.Run(":8080")
}
