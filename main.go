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

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.POST("/echo", func(c *gin.Context) {
		var req request
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	router.Run(":8080")

}
