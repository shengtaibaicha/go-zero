package models

import "time"

// Users 对应数据库中的users表
type Users struct {
	UserId         string    `gorm:"primary_key;column:user_id" json:"userId"`
	UserName       string    `gorm:"column:user_name" json:"userName"`
	UserEmail      string    `gorm:"column:user_email" json:"userEmail"`
	UserAvatar     string    `gorm:"column:user_avatar" json:"userAvatar"`
	UserPassword   string    `gorm:"column:user_password" json:"userPassword"`
	JoinDate       time.Time `gorm:"column:join_date;type:datetime" json:"joinDate"`
	DownloadNumber int64     `gorm:"column:download_number" json:"downloadNumber"`
	UploadNumber   int64     `gorm:"column:upload_number" json:"uploadNumber"`
	CollectNumber  int64     `gorm:"column:collect_number" json:"collectNumber"`
	Role           string    `gorm:"column:role" json:"role"`
	Enable         int       `gorm:"column:enable" json:"enable"`
}

// TableName 指定表名
func (Users) TableName() string {
	return "users" // 你的表名
}
