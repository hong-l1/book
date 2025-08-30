package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/events/article"
)

type App struct {
	Server    *gin.Engine
	Consumers []article.Consumer
}
