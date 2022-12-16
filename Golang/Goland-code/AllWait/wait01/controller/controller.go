package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wait01/dao"
	"wait01/model"
)

//一、用户操作
//注册--增加一个用户
//登录--查找一个用户
//修改用户
//显示--显示所有用户
//注销--删除一个用户

// 注册--增加一个用户
func RegisterUser(c *gin.Context) {
	//接收 发送的 用户名 和 密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 判断用户是否存在
	var status string //状态
	var err error
	count := -1
	u := dao.Mgr.Login(username)
	if u.Username == "" && password != "" && username != "" {
		//导入user结构
		_, count = dao.Mgr.GetAllUser() //已有总编号数
		user := model.User{
			ID:       count + 1, //在已有总编号数上增加
			Username: username,
			Password: password,
		}
		//正式创建用户
		err = dao.Mgr.RegisterUser(&user)
		if err == nil {
			status = "注册成功，请登录"
			fmt.Printf("username : %s ; password : %s \n", username, password)
			//c.Redirect(301, "/login")	//跳转到登录页面
		} else {
			status = "错误，注册失败"
		}
	} else {
		if password == "" || username == "" {
			status = "错误，用户或密码不能为空"
		} else {
			status = "错误，用户已存在"
		}
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":   status,
		"error":    err,
		"ID":       count + 1,
		"Username": username,
		"Password": password,
	})
	//获取所有用户
	GetAllUser(c)
}

// 登录--查找一个用户
func Login(c *gin.Context) {
	//接收 发送的 用户名 和 密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	var status string //状态
	// 判断用户是否存在
	u := dao.Mgr.Login(username)
	if u.Username == "" {
		status = "错误，用户名不存在"
	} else {
		if u.Password != password {
			status = "密码错误"
		} else {
			status = "登陆成功"
			fmt.Printf("ID : %d Username : %s   Password : %s \n", u.ID, u.Username, u.Password)
		}
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":   status,
		"ID":       u.ID,
		"Username": u.Username,
		"Password": u.Password,
	})
}

// 修改用户
func UpdateUser(c *gin.Context) {
	//用户名，新密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	//先查找
	u := dao.Mgr.Login(username)
	var status string
	var err error
	//再判断
	if u.Username != "" && password != "" {
		//映射，便于更新
		t := make(map[string]interface{})
		t["password"] = password
		//实现更新
		err = dao.Mgr.UpdateUser(u, t)
		if err == nil {
			status = "更新成功"
		} else {
			status = "更新失败"
		}
	} else {
		if password == "" {
			status = "错误，密码不能为空"
		} else {
			status = "错误，用户不存在"
		}
	}
	c.JSON(200, gin.H{
		"status":   status,
		"error":    err,
		"ID":       u.ID,
		"Username": u.Username,
		"password": u.Password, //旧密码
		"Password": password,   //新密码
	})
	//获取所有用户
	GetAllUser(c)
}

// 显示--显示所有用户
func GetAllUser(c *gin.Context) {
	//获取所有用户切片和用户数
	users, count := dao.Mgr.GetAllUser()
	//利用map集合来排序
	u := make(map[int]model.User)
	//以ID为key，用户结构体为value
	for user := range users {
		u[users[user].ID] = users[user]
	}
	//按Id由小到大的顺序(也是创建的时间先后顺序)，显示用户名
	var i int
	for i = 1; i <= count; i++ {
		c.JSON(200, gin.H{
			"ID":       u[i].ID,
			"Username": u[i].Username,
			"Password": u[i].Password,
		})
		fmt.Printf("ID : %d Username : %s   Password : %s \n", u[i].ID, u[i].Username, u[i].Password)
	}
	fmt.Printf("已注册用户数为：%d\n", count)
}

// 注销--删除一个用户
func DeleteUser(c *gin.Context) {
	//接收 发送的 用户名
	username := c.PostForm("username")
	//先查找
	u := dao.Mgr.Login(username)
	var begin int
	var status string
	var err bool
	//判断
	if u.Username != "" {
		//正式删除用户
		err, begin = dao.Mgr.DeleteUser(u)
		if !err {
			status = "删除成功"
			//获取所有用户切片和用户数
			users, count := dao.Mgr.GetAllUser()
			//利用map集合来排序
			ur := make(map[int]model.User)
			//以ID为key，用户结构体为value
			for user := range users {
				ur[users[user].ID] = users[user]
			}
			//修改变动ID
			var i int
			t := make(map[string]interface{})
			for i = begin; i <= count+1; i++ {
				t["ID"] = i - 1
				dao.Mgr.UpdateUser(ur[i], t)
			}
		} else {
			status = "删除失败"
		}
	} else {
		status = "错误，用户不存在"
	}
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":   status,
		"ID":       begin,
		"Username": username,
	})
	//获取所有用户
	GetAllUser(c)
}

//二、等待操作
//显示所有等待
//---下面要先登录---//
//查找一个等待
//修改一个等待
//增加一个等待
//删除一个等待

// 显示所有等待
func GetAllWait(c *gin.Context) {
	//获取所有等待切片和等待数
	waits, count := dao.Mgr.GetAllWait()
	//利用map集合来排序
	w := make(map[int]model.Wait)
	//以ID为key，等待结构体为value
	for wait := range waits {
		w[waits[wait].ID] = waits[wait]
	}
	//按Id由小到大的顺序(也是创建的时间先后顺序)，显示等待标题
	var i int
	for i = 1; i <= count; i++ {
		c.JSON(200, gin.H{
			"ID":        w[i].ID,
			"Title":     w[i].Title,
			"Number":    w[i].Number,
			"Address":   w[i].Address,
			"Transport": w[i].Transport,
		})
		fmt.Printf("ID : %d \nTitle : %s \nNumber : %s \nAddress : %s \nTransport : %s\n\n", w[i].ID, w[i].Title, w[i].Number, w[i].Address, w[i].Transport)
	}
	fmt.Printf("记录等待数为：%d\n", count)
}

// 查找一个等待
func GetWait(c *gin.Context) {
	//接收 发送的 等待
	title := c.PostForm("title") //主题

	var status string //状态
	//判断等待是否存在
	w := dao.Mgr.GetWait(title)
	if w.Title == "" {
		status = "错误，等待不存在"
	} else {
		status = "查找成功"
		fmt.Printf("ID : %d \nTitle : %s \nNumber : %s \nAddress : %s \nTransport : %s\n\n", w.ID, w.Title, w.Number, w.Address, w.Transport)
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":    status,
		"ID":        w.ID,
		"Title":     w.Title,
		"Number":    w.Number,
		"Address":   w.Address,
		"Transport": w.Transport,
	})
}

// 修改一个等待
func UpdateWait(c *gin.Context) {
	//等待主题，新的内容
	title := c.PostForm("title")         //主题
	number := c.PostForm("number")       //人数
	address := c.PostForm("address")     //地址
	transport := c.PostForm("transport") //前往方式
	var status string                    //状态
	var err error
	//判断等待是否存在
	w := dao.Mgr.GetWait(title)
	if w.Title != "" {
		//映射，便于更新
		t := make(map[string]interface{})
		t["number"] = number
		t["address"] = address
		t["transport"] = transport
		//删除空值
		for wi, wj := range t {
			if wj == "" {
				delete(t, wi)
			}
		}
		//实现更新
		err = dao.Mgr.UpdateWait(w, t)
		if err == nil {
			status = "更新成功"
		} else {
			status = "更新失败"
		}
	} else {
		status = "错误，等待不存在"
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":    status,
		"error":     err,
		"ID":        w.ID,
		"Title":     w.Title,
		"number":    w.Number,
		"Number":    number,
		"address":   w.Address,
		"Address":   address,
		"transport": w.Transport,
		"Transport": transport,
	})
	//获取所有等待
	GetAllWait(c)
}

// 增加一个等待
func AddWait(c *gin.Context) {
	//接收等待
	title := c.PostForm("title")         //主题
	number := c.PostForm("number")       //人数
	address := c.PostForm("address")     //地址
	transport := c.PostForm("transport") //前往方式
	//判断等待是否存在
	count := -1
	var status string //状态
	var err error
	w := dao.Mgr.GetWait(title)
	if w.Title == "" {
		//导入wait结构
		_, count = dao.Mgr.GetAllWait() //已有总编号数
		wait := model.Wait{
			ID:        count + 1, //在已有总编号数上增加
			Title:     title,
			Number:    number,
			Address:   address,
			Transport: transport,
		}
		//正式添加等待
		err = dao.Mgr.AddWait(&wait)
		if err == nil {
			status = "添加成功"
			fmt.Printf("ID : %d \nTitle : %s \nNumber : %s \nAddress : %s \nTransport : %s\n\n", w.ID, w.Title, w.Number, w.Address, w.Transport)
		} else {
			status = "错误，添加失败"
		}
	} else {
		status = "错误，等待已存在"
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":    status,
		"ID":        count + 1,
		"Title":     title,
		"Number":    number,
		"Address":   address,
		"Transport": transport,
	})
	//获取所有等待
	GetAllWait(c)
}

// 删除一个等待
func DeleteWait(c *gin.Context) {
	//接收 发送的 等待
	title := c.PostForm("title") //主题
	//判断等待是否存在
	w := dao.Mgr.GetWait(title)
	var status string //状态
	var begin int
	var err bool
	if w.Title != "" {
		//正式删除
		err, begin = dao.Mgr.DeleteWait(w)
		if !err {
			status = "删除成功"
			//获取删除后所有等待切片和等待数
			waits, count := dao.Mgr.GetAllWait()
			//利用map集合来排序
			wr := make(map[int]model.Wait)
			//以ID为key，等待结构体为value
			for wait := range waits {
				wr[waits[wait].ID] = waits[wait]
			}
			//修改变动ID
			var i int
			t := make(map[string]interface{})
			for i = begin; i <= count+1; i++ {
				t["ID"] = i - 1
				dao.Mgr.UpdateWait(wr[i], t)
			}
		} else {
			status = "删除失败"
		}
	} else {
		status = "错误，等待不存在"
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":    status,
		"ID":        w.ID,
		"Title":     w.Title,
		"Number":    w.Number,
		"Address":   w.Address,
		"Transport": w.Transport,
	})
	//获取所有等待
	GetAllWait(c)
}
