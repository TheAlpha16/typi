package models

type Video struct {
	VID         string `gorm:"primaryKey;column:vid" json:"vid"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	PublishedAt string `gorm:"column:published_at" json:"published_at"`
	Thumbnail   string `gorm:"column:thumbnail" json:"thumbnail"`
}

type Config struct {
	Key   string `gorm:"primaryKey;column:key"`
	Value string `gorm:"column:value"`
}

func (Config) TableName() string {
	return "config"
}
