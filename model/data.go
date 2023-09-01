package model

import "time"

type Data struct {
	Id        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	Count     uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Data) TableName() string {
	return "data"
}
