package controller

import (
	"fmt"
	"log"
	"strconv"
	"wait/dao"
	"wait/model"

	"github.com/gin-gonic/gin"
)

// 用户注册

func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := model.User{
		Username: username,
		Password: password,
	}

	dao.Mgr.RegisterUser(&user)

	c.Redirect(301, "/")
}

// 实现用户注册页面

func GoRegister(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

// 用户登录

func Login(c *gin.Context) {
	//接收 发送的 用户名 和 密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println(username)

	u := dao.Mgr.Login(username) //对用户名判断

	if u.Username == "" {
		c.HTML(200, "login.html", "用户名不存在")
		fmt.Println("用户名不存在")
	} else {
		if u.Password != password {
			c.HTML(200, "login.html", "密码错误")
			fmt.Println("密码错误")
		} else {
			c.Redirect(301, "/")
			fmt.Println("登陆成功")
		}
	}
}

// 实现用户登陆页面

func GoLogin(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

// 首页

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

// 列表

func ListUser(c *gin.Context) {
	c.HTML(200, "userlist.html", nil)
}

// 操作等待
// 获取等待列表

func GetWaitIndex(c *gin.Context) {
	waits := dao.Mgr.GetAllWait()
	c.HTML(200, "postIndex.html", waits)
}

// 添加等待		//  需 要 大 改

func AddWait(c *gin.Context) {
	title := c.PostForm("title")
	number := c.PostForm("number")
	address := c.PostForm("address")
	transport := c.PostForm("transport")
	wait := model.Wait{
		Title:     title,
		Number:    number,
		Address:   address,
		Transport: transport,
	}

	dao.Mgr.AddWait(&wait)

	c.Redirect(302, "/post_index")
}

func Quitwait(c *gin.Context) {
	c.HTML(200, "quit.html", nil)
	//c.Redirect(303, "https://cn.bing.com/")

}

// 跳转到添加等待

func GoAddWait(c *gin.Context) {
	c.HTML(200, "post.html", nil)
}

// 显示等待详细

func WaitDetail(c *gin.Context) {
	s := c.Query("pid")

	pid, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}

	p := dao.Mgr.GetWait(pid)

	c.HTML(200, "detail.html", gin.H{
		"Title":     p.Title,
		"Number":    p.Number,
		"Address":   p.Address,
		"Transport": p.Transport,
	})
}
