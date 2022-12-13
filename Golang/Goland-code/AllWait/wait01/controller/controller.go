package controller

import (
	"github.com/gin-gonic/gin"
)

//获取首页界面

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
