package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CORRELATION_HEADER = "Correlation-Id"
)

type request struct {
	Name   string         `json:"name" binding:"required"`
	Nested []nestedStruct `json:"nested" binding:"required,dive"`
}

type nestedStruct struct {
	Other  string  `json:"other" binding:"required"`
	Number float64 `json:"number" binding:"required,gte=0"`
}

type logMsg struct {
	ClientIP      string
	TimeStamp     string
	Method        string
	Path          string
	Status        int
	Latency       string
	ErrorMessage  string
	CorrelationID string
}

func echoHandler(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func livenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func logFormatter(param gin.LogFormatterParams) string {

	correlationID := ""
	ch := param.Request.Header[CORRELATION_HEADER]
	if len(ch) > 0 {
		correlationID = ch[0]
	}

	m := logMsg{
		ClientIP:      param.ClientIP,
		TimeStamp:     param.TimeStamp.Format(time.RFC1123),
		Method:        param.Method,
		Path:          param.Path,
		Status:        param.StatusCode,
		Latency:       param.Latency.String(),
		ErrorMessage:  param.ErrorMessage,
		CorrelationID: correlationID,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "error formatting log:" + fmt.Sprintf("%+v\n", m) + "\n"
	}
	return string(b) + "\n"

}

func main() {

	gin.SetMode(gin.ReleaseMode)

	echoRouter := gin.New()
	echoRouter.Use(gin.Recovery())
	echoRouter.SetTrustedProxies(nil)
	echoRouter.Use(gin.LoggerWithFormatter(logFormatter))
	echoRouter.POST("/echo", echoHandler)

	livenessRouter := gin.New()
	livenessRouter.Use(gin.Recovery())
	livenessRouter.SetTrustedProxies(nil)
	livenessRouter.Use(gin.LoggerWithFormatter(logFormatter))
	livenessRouter.GET("/liveness", livenessHandler)

	go echoRouter.Run(":8080")
	livenessRouter.Run(":8081")
}
