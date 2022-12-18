package model

//model与数据库中的实体一一对应
import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int    `json:"id"`       //编号
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}

type Wait struct {
	gorm.Model
	ID        int    `json:"id"`        //编号
	Title     string `json:"title"`     //等待主题
	Number    string `json:"number"`    //等待人数
	Wtime     int    `json:"wtime"`     //等待时间
	Address   string `json:"address"`   //等待地址
	Transport string `json:"transport"` //前往方式
	Imgindex  int    `json:"imgindex"`  //图片循环下标
	Sum       int    `json:"sum"`       //对应存储图片数
}

type Img struct {
	gorm.Model
	Title  string `json:"title"`  //等待主题
	Name   string `json:"name"`   //图片名
	Index  int    `json:"index"`  //下标
	Data   string `json:"data"`   //图片内容
	Imgurl string `json:"imgurl"` //图片路径
}
