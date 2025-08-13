package models

import (
	"gorm.io/plugin/soft_delete"
)

type Files struct {
	FileId     string `gorm:"primaryKey;column:file_id" json:"fileId"`
	FileUrl    string `gorm:"column:file_url" json:"fileUrl"`
	UploadTime string `gorm:"column:upload_time" json:"uploadTime"`
	Status     string `gorm:"column:status" json:"status"`
	UserId     string `gorm:"column:user_id" json:"userId"`
	FileName   string `gorm:"column:file_name" json:"fileName"`
	Number     int64  `gorm:"column:number" json:"number"`
	FileTitle  string `gorm:"column:file_title" json:"fileTitle"`
	FileUrlse  string `gorm:"column:file_urlse" json:"fileUrlse"`

	// 这个gorm插件不加softDelete:flag时使用unix时间戳作为删除标志，加上则以0,1作为删除标志
	// gorm.DeletedAt 这个类型使用*time.Time作为删除标志，其他类型需要使用插件实现

	//Deleted gorm.DeletedAt `gorm:"column:deleted;softDelete:flag" json:"deleted"`
	Deleted soft_delete.DeletedAt `gorm:"column:deleted;softDelete:flag" json:"deleted"`
}

func (Files) TableName() string {
	return "files" // 你的表名
}
