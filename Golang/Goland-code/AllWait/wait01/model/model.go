package model

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
	Imgurl    string `json:"imgurl"`    //图片地址
}
