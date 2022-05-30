package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikeallman/rest-api/log"
)

type request struct {
	Name   string         `json:"name" binding:"required"`
	Nested []nestedStruct `json:"nested" binding:"required,dive"`
}

type nestedStruct struct {
	Other  string  `json:"other" binding:"required"`
	Number float64 `json:"number" binding:"required,gte=0"`
}

func echoPOSTHandler(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func livenessGETHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func baseRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)
	r.Use(gin.LoggerWithFormatter(log.Formatter))
	return r
}

func echoRouter() *gin.Engine {
	r := baseRouter()
	r.POST("/echo", echoPOSTHandler)
	return r
}

func livenessRouter() *gin.Engine {
	r := baseRouter()
	r.GET("/liveness", livenessGETHandler)
	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	echoRouter := echoRouter()
	livenessRouter := livenessRouter()

	go echoRouter.Run(":8080")
	livenessRouter.Run(":8081")
}
