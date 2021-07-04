package controller

import (
	"github.com/gin-gonic/gin"
	"golang-student/model"
	"net/http"
)

// LoginStudent 登录系统
func LoginStudent(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "账号或密码不能为空",
		})
		return
	}
	data, err := model.LoginStudent(username, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": &data,
			"msg":  "success",
		})
	}

}

func GetScore(c *gin.Context) {
	semesterId := c.Query("semester_id")
	token := c.Query("token")

	data, err := model.GetScore(token, semesterId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": &data,
			"msg":  "success",
		})
	}

}
