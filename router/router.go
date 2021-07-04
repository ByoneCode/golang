package router

import (
	"github.com/gin-gonic/gin"
	"golang-student/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(Cors())
	// 登录系统
	r.POST("/login", controller.LoginStudent)
	r.GET("/score", controller.GetScore)

	return r
}
