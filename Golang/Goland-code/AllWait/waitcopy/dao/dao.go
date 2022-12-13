package dao

import (
	"log"
	"wait/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Manager interface {
	RegisterUser(user *model.User)
	Login(username string) model.User

	// 等 待 操 作

	AddWait(wait *model.Wait)   // 添加等待
	GetAllWait() []model.Wait   // 获取多个等待
	GetWait(pid int) model.Wait // 获取单个等待
}

type manager struct {
	db *gorm.DB
}

var Mgr Manager

func init() {
	dsn := "root:20050901@mysqlHrx@(127.0.0.1:3306)/wait?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}

	Mgr = &manager{db: db}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Wait{})
}

func (mgr *manager) RegisterUser(user *model.User) {
	mgr.db.Create(user)
}

//对登录的用户名进行判断

func (mgr *manager) Login(username string) model.User {
	var user model.User
	//
	mgr.db.Where("username=?", username).First(&user)
	return user
}

// 等待操作
func (mgr *manager) AddWait(wait *model.Wait) {
	mgr.db.Create(wait)
}

func (mgr *manager) GetAllWait() []model.Wait {
	var waits = make([]model.Wait, 1000)
	mgr.db.Find(&waits)
	return waits
}

func (mgr *manager) GetWait(pid int) model.Wait {
	var wait model.Wait
	mgr.db.First(&wait, pid)
	return wait
}
