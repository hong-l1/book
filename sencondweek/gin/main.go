package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "this is post")
	})
	server.GET("/get/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, name)
	})
	server.GET("/view/*.html", func(c *gin.Context) {
		name := c.Param(".html")
		c.String(http.StatusOK, name)
	})
	server.GET("/order", func(c *gin.Context) {
		id := c.Query("id")
		c.String(http.StatusOK, "查询参数"+id)
	})
	server.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
