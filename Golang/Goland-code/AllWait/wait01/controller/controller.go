package controller

//控制层，负责具体模块的业务流程控制

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
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
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
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
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
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
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
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
	username, _ := c.GetPostForm("username")
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
			"Wtime":     w[i].Wtime,
			"Imgindex":  w[i].Imgindex,
			"Sum":       w[i].Sum,
			"Address":   w[i].Address,
			"Transport": w[i].Transport,
		})
		fmt.Printf("ID : %d \nTitle : %s \nNumber : %s\nWtime : %d \nImgindex : %d \nSum : %d\nAddress : %s \nTransport : %s\n\n",
			w[i].ID, w[i].Title, w[i].Number, w[i].Wtime, w[i].Imgindex, w[i].Sum, w[i].Address, w[i].Transport)
	}
	fmt.Printf("记录等待数为：%d\n", count)
}

// 查找一个等待
func GetWait(c *gin.Context) {
	//接收 发送的 等待
	title, _ := c.GetPostForm("title") //主题

	var status string //状态
	//判断等待是否存在
	w := dao.Mgr.GetWait(title)
	if w.Title == "" {
		status = "错误，等待不存在"
	} else {
		status = "查找成功"
		fmt.Printf("ID : %d \nTitle : %s \nNumber : %s\nWtime : %d \nImgindex :%d \nSum : %d \nAddress : %s \nTransport : %s\n\n",
			w.ID, w.Title, w.Number, w.Wtime, w.Imgindex, w.Sum, w.Address, w.Transport)
	}
	//响应
	fmt.Printf("%s\n", status)
	c.JSON(200, gin.H{
		"status":    status,
		"ID":        w.ID,
		"Title":     w.Title,
		"Wtime":     w.Wtime,
		"Imgindex":  w.Imgindex,
		"Sum":       w.Sum,
		"Number":    w.Number,
		"Address":   w.Address,
		"Transport": w.Transport,
	})
	GetImgs(c)
}

// 修改一个等待
func UpdateWait(c *gin.Context) {
	//等待主题，新的内容
	title, _ := c.GetPostForm("title")         //主题
	number, _ := c.GetPostForm("number")       //人数
	address, _ := c.GetPostForm("address")     //地址
	transport, _ := c.GetPostForm("transport") //前往方式
	var status string                          //状态
	var err error
	var wtime int
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
		wtime, err = strconv.Atoi(number)
		if err == nil && wtime >= 0 {
			wtime = ((wtime-1)/5 + 1) * 5
			t["wtime"] = wtime
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
		"wtime":     w.Wtime,
		"Wtime":     wtime,
		"Imgindex":  w.Imgindex,
		"Sum":       w.Sum,
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
	title, _ := c.GetPostForm("title")         //主题
	number, _ := c.GetPostForm("number")       //人数
	address, _ := c.GetPostForm("address")     //地址
	transport, _ := c.GetPostForm("transport") //前往方式
	//判断等待是否存在
	count := -1
	var status string //状态
	var err error
	var wtime int
	w := dao.Mgr.GetWait(title)
	if w.Title == "" {
		//导入wait结构
		_, count = dao.Mgr.GetAllWait() //已有总编号数
		wtime, err = strconv.Atoi(number)
		if err == nil && wtime >= 0 {
			wtime = ((wtime-1)/5 + 1) * 5
			wait := model.Wait{
				ID:        count + 1, //在已有总编号数上增加
				Title:     title,
				Wtime:     wtime,
				Number:    number,
				Address:   address,
				Transport: transport,
			}
			//正式添加等待
			err = dao.Mgr.AddWait(&wait)
			if err == nil {
				status = "添加成功"
				//上传图片
				Upload(c)
				fmt.Printf("ID : %d \nTitle : %s \nNumber : %s \nWtime : %d \nAddress : %s \nTransport : %s\n",
					count+1, title, number, wtime, address, transport)
			} else {
				status = "错误，添加失败"
			}
		} else {
			status = "错误，非法人数"
		}

	} else {
		status = "错误，等待已存在"
	}
	//响应
	fmt.Printf("%s\n\n", status)
	c.JSON(200, gin.H{
		"status":    status,
		"Error":     err,
		"ID":        count + 1,
		"Title":     title,
		"Number":    number,
		"Wtime":     wtime,
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
		//删除图片
		dao.Mgr.DeleteImgs(w.Title)
		os.RemoveAll("./upload/" + w.Title)
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
		"Wtime":     w.Wtime,
		"Imgindex":  w.Imgindex,
		"Sum":       w.Sum,
		"Address":   w.Address,
		"Transport": w.Transport,
	})
	//获取所有等待
	GetAllWait(c)
}

//图片操作
// 图片上传
// 同一title下图片的获取
// 图片删除

// 同一title下图片的获取
func GetImgs(c *gin.Context) {
	title, _ := c.GetPostForm("title")
	w := dao.Mgr.GetWait(title)
	is := dao.Mgr.GetImgs(title)
	os.Mkdir("./upload/"+title, 0777)
	t := make(map[int]model.Img)
	for _, v := range is {
		t[v.Index] = v
	}
	for k := 1; k <= w.Imgindex; k++ {
		c.JSON(200, gin.H{
			"Index":  t[k].Index,
			"Name":   t[k].Name,
			"Imgurl": t[k].Imgurl,
		})
		fmt.Printf("Index : %d\nName : %s\nImgurl :%s\n",
			t[k].Index, t[k].Name, t[k].Imgurl)
		ddd, _ := base64.StdEncoding.DecodeString(t[k].Data)
		os.WriteFile(t[k].Imgurl, ddd, 0666)
	}
}

// 图片上传
func Upload(c *gin.Context) {
	var status string
	//获取文件，icon实现对上传文件的访问，header是对上传文件信息的标记
	icon, header, err := c.Request.FormFile("file")
	title, _ := c.GetPostForm("title")
	//判断是否已有同名图片
	i := dao.Mgr.GetImg(title, header.Filename)
	if i.Name != "" {
		status = "错误，该主题下已存在这一图片名"
	} else {
		if err == nil {
			defer icon.Close()
			//path.Ext是取后缀，Tolower小写
			//ext := strings.ToLower(path.Ext(header.Filename))
			if header.Size > 1024*1024*2 {
				fmt.Println("文件过大")
			}
			buf := bytes.NewBuffer(nil)
			//读取icon的数据存入buf中
			if _, err1 := io.Copy(buf, icon); err1 != nil {
				return
			}
			//转码成 base64 来存储
			data := base64.StdEncoding.EncodeToString(buf.Bytes())
			//找到对应的 等待
			w := dao.Mgr.GetWait(title)
			sum := w.Sum + 1          //等待中的总图片数
			index := w.Imgindex%5 + 1 //指向当前图片的下标//实现循环
			if sum == 6 {             //超过 5 ，总数不变
				sum = 5
				is := dao.Mgr.GetImgs(title) //若已满 ， 则会先删除 再存储
				t := make(map[int]model.Img)
				for _, v := range is {
					t[v.Index] = v
				}
				err2 := dao.Mgr.DeleteImg(t[index])
				os.Remove(t[index].Imgurl)
				if err2 != nil {
					fmt.Println("错误，删除失败\n")
					c.JSON(200, gin.H{
						"status": "错误，删除失败",
					})
				}
			}
			//转换成字符串，存在 url 里
			sindex := strconv.Itoa(index)
			imgurl := "./upload/" + title + "/" + sindex + "-" + header.Filename
			//导入 Img 结构
			img := model.Img{
				Title:  title,
				Index:  index,
				Name:   header.Filename,
				Data:   data,
				Imgurl: imgurl,
			}
			//正式上传
			err0 := dao.Mgr.Upload(img)
			if err0 == nil {
				status = "上传成功"
				fmt.Printf("Title : %s\nIndex : %d\nName : %s\nImgurl :%s\n",
					title, index, header.Filename, imgurl)
				//更新 等待 中的图片数
				t := map[string]interface{}{
					"Imgindex": index, //指向当前图片的下标
					"Sum":      sum,   //等待中的总图片数
				}
				dao.Mgr.UpdateWait(w, t)
			} else {
				status = "上传失败"
			}
			//ntime := time.Now().Format("20060102150405")
			//没有，则创建文件目录；有，无影响
			os.Mkdir("./upload/"+title, 0777)
			ddd, _ := base64.StdEncoding.DecodeString(data)
			os.WriteFile(imgurl, ddd, 0666)
			c.JSON(200, gin.H{
				"Title":  title,
				"Index":  index,
				"Name":   header.Filename,
				"Imgurl": imgurl,
			})
		} else {
			status = "未上传图片"
		}
	}
	fmt.Println(status, "\n")
	c.JSON(200, gin.H{
		"status": status,
	})
}

/*
func Upload(c *gin.Context) {
	fmt.Printf("OK!\n")

	fileheader, _ := c.FormFile("file")
	name := fileheader.Filename
	size := fileheader.Size
	header := fileheader.Header
	fmt.Printf("name[%s],size[%d],header[%#v]\n", name, size, header)
	c.JSON(200, gin.H{
		"name":   name,
		"size":   size,
		"header": header,
	})

	extName := path.Ext(name)
	allowExtMap := map[string]bool{
		".jpg": true,
		".png": true,
	}
	var status string

	if _, ok := allowExtMap[extName]; !ok {
		status = "上传错误，文件类型不合法"
	} else {
		time := time.Now().Format("20060102150405")
		savedir := path.Join("/upload/", time+extName)
		err := c.SaveUploadedFile(fileheader, savedir)
		if err != nil {
			status = "文件保存失败"
		} else {
			status = "上传成功"
			imgurl := strings.Replace(savedir, "/", "/", -1)
			c.JSON(200, gin.H{
				"imgurl": imgurl,
			})
		}
	}

	fmt.Printf("%s\n", status)
}
*/
