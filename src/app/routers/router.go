package routers

import (
	v1 "github.com/devops-salt/src/app/routers/api/v1"
	"github.com/devops-salt/src/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(config.GetRunMod())

	r.GET("/task/download/:addr/:task_id", v1.DownTask)

	r.POST("/task/add", v1.AddTask)

	r.POST("/callback", v1.Callback)

	return r
}


