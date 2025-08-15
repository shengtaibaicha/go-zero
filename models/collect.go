package models

type Collect struct {
	Id     int32  `gorm:"primary_key;column:id;autoIncrement" json:"id"`
	FileId string `gorm:"column:file_id" json:"fileId"`
	UserId string `gorm:"column:user_id" json:"userId"`
}

func (Collect) TableName() string {
	return "collect" // 你的表名
}
