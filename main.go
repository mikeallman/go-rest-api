package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	Name   string         `json:"name" binding:"required"`
	Nested []nestedStruct `json:"nested" binding:"required,dive"`
}

type nestedStruct struct {
	Other  string  `json:"other" binding:"required"`
	Number float64 `json:"number" binding:"required,gte=0"`
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

func main() {

	gin.SetMode(gin.ReleaseMode)

	echoRouter := gin.Default()
	echoRouter.SetTrustedProxies(nil)
	echoRouter.POST("/echo", echoHandler)

	livenessRouter := gin.Default()
	livenessRouter.SetTrustedProxies(nil)
	livenessRouter.GET("/liveness", livenessHandler)

	go echoRouter.Run(":8080")
	livenessRouter.Run(":8081")
}
