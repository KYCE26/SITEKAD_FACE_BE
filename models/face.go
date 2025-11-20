package models

import (
	"encoding/json"
	"time"
)

type UserFace struct {
	Id           int64           `gorm:"primaryKey" json:"id"`
	PenempatanId int64           `gorm:"index" json:"penempatan_id"`
	Name         string          `json:"name"`
	Embedding    json.RawMessage `gorm:"type:json" json:"-"`        // Raw JSON dari DB
	Vector       []float64       `gorm:"-" json:"embedding"`        // Helper buat coding
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
}

func (UserFace) TableName() string {
	return "user_faces"
}