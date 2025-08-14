package models

import "time"

// Users 对应数据库中的users表
type Users struct {
	UserId         string    `gorm:"primary_key;column:user_id"`
	UserName       string    `gorm:"column:user_name"`
	UserEmail      string    `gorm:"column:user_email"`
	UserAvatar     string    `gorm:"column:user_avatar"`
	UserPassword   string    `gorm:"column:user_password"`
	JoinDate       time.Time `gorm:"column:join_date;type:datetime"`
	DownloadNumber int64     `gorm:"column:download_number"`
	UploadNumber   int64     `gorm:"column:upload_number"`
	CollectNumber  int64     `gorm:"column:collect_number"`
	Role           string    `gorm:"column:role"`
	Enable         int       `gorm:"column:enable"`
}

// TableName 指定表名
func (Users) TableName() string {
	return "users" // 你的表名
}
