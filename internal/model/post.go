package model

import "time"

type Post struct {
	Id          int       `gorm:"primaryKey;type:bigint"`
	Title       string    `gorm:"type:varchar(255)"`
	Content     string    `gorm:"type:text"`
	LikeCount   int       `gorm:"type:int"`
	PublishTime time.Time `gorm:"type:datetime"`
	PublisherId int       `gorm:"type:bigint"`
}
