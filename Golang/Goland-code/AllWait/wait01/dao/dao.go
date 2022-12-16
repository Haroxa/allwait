package dao

import (
	"log"
	"wait01/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//操作管理 的接口

type Manager interface {

	//一、用户操作
	//注册--增加一个用户
	//登录--查找一个用户
	//修改用户
	//显示--显示所有用户
	//注销--删除一个用户

	RegisterUser(user *model.User) error
	Login(username string) model.User
	UpdateUser(user model.User, t map[string]interface{}) error
	GetAllUser() /*[]model.User */ ([]model.User, int)
	DeleteUser(user model.User) (bool, int)

	//二、等待操作
	//显示所有等待
	//---下面要先登录---//
	//查找一个等待
	//修改一个等待
	//增加一个等待
	//删除一个等待

	GetAllWait() ([]model.Wait, int) // 获取多个等待
	GetWait(title string) model.Wait // 获取单个等待
	UpdateWait(wait model.Wait, t map[string]interface{}) error
	AddWait(wait *model.Wait) error // 添加等待
	DeleteWait(wait model.Wait) (bool, int)
}

type manager struct {
	db *gorm.DB
}

var Mgr Manager

//初始化，连接到数据库

func init() {
	dsn := "root:20050901@mysqlHrx@(127.0.0.1:3306)/wait02?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}

	Mgr = &manager{db: db}
	db.AutoMigrate(&model.User{})

	db.AutoMigrate(&model.Wait{})
}

// 接口内 函数 定义

//一、用户操作
//注册--增加一个用户
//登录--查找一个用户
//修改用户
//显示--显示所有用户
//注销--删除一个用户

// 注册--增加一个用户
func (mgr *manager) RegisterUser(user *model.User) error {
	//在数据库中创建一个新用户
	return mgr.db.Create(user).Error
}

// 登录--查找一个用户
func (mgr *manager) Login(username string) model.User {
	var user model.User
	//在数据库中查找用户名，找到就给user赋值,找不到就为空
	//正常注册时会出现 record not found
	//First() 函数找不到record的时候，会返回error: record not found ，
	//而Find() 则是返回nil
	mgr.db.Where("username=?", username).First(&user)
	return user
}

// 修改用户
func (mgr *manager) UpdateUser(user model.User, t map[string]interface{}) error {
	//对指定的用户模型，更新多个字段值，并判断是否有误
	return mgr.db.Model(&user).Updates(t).Error
}

// 显示--显示所有用户
func (mgr *manager) GetAllUser() ([]model.User, int) {
	users := make([]model.User, 1000) //user模型切片
	var c int64
	mgr.db.Find(&users).Count(&c) //找到所有用户并赋值
	//Count函数，直接返回查询匹配的行数。
	return users, int(c)
}

// 注销--删除已登录用户
func (mgr *manager) DeleteUser(user model.User) (bool, int) {
	//删除指定用户以及数据库中的删除记录
	mgr.db.Unscoped().Delete(&user)
	return false, user.ID
}

//二、等待操作
//显示所有等待
//---下面要先登录---//
//查找一个等待
//修改一个等待
//增加一个等待
//删除一个等待

// 显示所有等待
func (mgr *manager) GetAllWait() ([]model.Wait, int) {
	waits := make([]model.Wait, 1000) //模型切片
	var c int64
	mgr.db.Find(&waits).Count(&c) //找到所有等待并赋值
	//Count函数，直接返回查询匹配的行数
	return waits, int(c)
}

// 查找一个等待
func (mgr *manager) GetWait(title string) model.Wait {
	var wait model.Wait
	//在数据库中查找等待主题，找到就给wait赋值,找不到就为空
	//正常注册时会出现 record not found
	//First() 函数找不到record的时候，会返回error: record not found ，
	//而Find() 则是返回nil
	mgr.db.Where("title=?", title).First(&wait)
	return wait
}

// 修改一个等待
func (mgr *manager) UpdateWait(wait model.Wait, t map[string]interface{}) error {
	//对指定的等待模型，更新多个字段值，并判断是否有误
	return mgr.db.Model(&wait).Updates(t).Error
}

// 增加一个等待
func (mgr *manager) AddWait(wait *model.Wait) error {
	//在数据库中创建一个新等待
	return mgr.db.Create(wait).Error
}

// 删除一个等待
func (mgr *manager) DeleteWait(wait model.Wait) (bool, int) {
	//删除指定用户以及数据库中的删除记录
	mgr.db.Unscoped().Delete(&wait)
	return false, wait.ID
}
