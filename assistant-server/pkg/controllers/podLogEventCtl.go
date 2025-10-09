package controllers

import (
	"context"
	"io"
	"net/http"

	"github.com/assistant-server/pkg/service"

	"github.com/gin-gonic/gin"
)

type PodLogEventCtl struct {
	podLogEventService *service.PodLogEventService
}

func NewPodLogEventCtl(service *service.PodLogEventService) *PodLogEventCtl {
	return &PodLogEventCtl{podLogEventService: service}
}

func (p *PodLogEventCtl) GetLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		ns := c.DefaultQuery("ns", "default")
		podName := c.DefaultQuery("podname", "")

		var tailLine int64 = 100
		req := p.podLogEventService.GetLogs(ns, podName, tailLine)

		rc, err := req.Stream(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer rc.Close()

		logData, err := io.ReadAll(rc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"data": string(logData)})
	}
}

func (p *PodLogEventCtl) GetEvent() func(c *gin.Context) {
	return func(c *gin.Context) {
		ns := c.DefaultQuery("ns", "default")
		podName := c.DefaultQuery("podname", "")

		e, err := p.podLogEventService.GetEvents(ns, podName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"data": e})
	}
}
