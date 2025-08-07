package models

type TagAndFile struct {
	Id     string `gorm:"primary_key column:id"`
	TagId  int32  `gorm:"column:tag_id"`
	FileId string `gorm:"column:file_id"`
}

func (TagAndFile) TableName() string {
	return "tagAndFile" // 你的表名
}
