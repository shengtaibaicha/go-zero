package models

type Files struct {
	FileId     string `gorm:"primary_key column:file_id" json:"fileId"`
	FileUrl    string `gorm:"column:file_url" json:"fileUrl"`
	UploadTime string `gorm:"column:upload_time" json:"uploadTime"`
	Status     string `gorm:"column:status" json:"status"`
	UserId     string `gorm:"column:user_id" json:"userId"`
	FileName   string `gorm:"column:file_name" json:"fileName"`
	Number     int64  `gorm:"column:number" json:"number"`
	FileTitle  string `gorm:"column:file_title" json:"fileTitle"`
	FileUrlse  string `gorm:"column:file_urlse" json:"fileUrlse"`
	Deleted    int    `gorm:"column:deleted;softDelete:1" json:"deleted"`
}

func (Files) TableName() string {
	return "files" // 你的表名
}
