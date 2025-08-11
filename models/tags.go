package models

type Tags struct {
	TagId   int32  `gorm:"primary_key column:tag_id" json:"tagId"`
	TagName string `gorm:"column:tag_name" json:"tagName"`
}

func (Tags) TableName() string {
	return "tags" // 你的表名
}
