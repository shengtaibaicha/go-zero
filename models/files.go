package models

type Files struct {
	FileId     string `gorm:"primary_key column:file_id"`
	FileUrl    string `gorm:"column:file_url"`
	UploadTime string `gorm:"column:upload_time"`
	Status     string `gorm:"column:status"`
	UserId     string `gorm:"column:user_id"`
	FileName   string `gorm:"column:file_name"`
	Number     int64  `gorm:"column:number"`
	FileTitle  string `gorm:"column:file_title"`
	FileUrlse  string `gorm:"column:file_urlse"`
}

func (Files) TableName() string {
	return "files" // 你的表名
}
