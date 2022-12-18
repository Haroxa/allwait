package router

import (
	"github.com/gin-gonic/gin"
	"wait01/controller"
)

func Start() {
	//创建服务
	e := gin.Default()

	//一、用户操作
	//注册--增加一个用户
	//登录--查找一个用户
	//修改用户
	//显示--显示所有用户
	//注销--删除一个用户
	e.POST("/register", controller.RegisterUser)
	e.POST("/login", controller.Login)
	e.POST("updateuser", controller.UpdateUser)
	e.GET("/getalluser", controller.GetAllUser)
	e.POST("/deleteuser", controller.DeleteUser)

	//二、等待操作
	//显示所有等待
	//---下面要先登录---//
	//查找一个等待
	//修改一个等待
	//增加一个等待
	//删除一个等待

	e.GET("/getallwait", controller.GetAllWait)
	e.POST("/getwait", controller.GetWait)
	e.POST("/updatewait", controller.UpdateWait)
	e.POST("/addwait", controller.AddWait)
	e.POST("/deletewait", controller.DeleteWait)

	e.POST("/upload", controller.Upload)
	//服务器端口，默认8080
	e.Run(":8080")
}
