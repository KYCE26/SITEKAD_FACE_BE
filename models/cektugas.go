package models

import "time"

type CekTugas struct {
	Ctid      int64 `gorm:"primaryKey" json:"ctid"`
	Ptid      int64 `gorm:"type bigint" json:"ptid"`
	Cid       int64 `gorm:"type bigint" json:"cid"`
	WaktuScan string `gorm:"type timestamp" json:"waktu_scan"`
	CreatedAt time.Time `gorm:"type timestamp" json:"created_at"`
}

func (CekTugas) TableName() string {
	return "cek_tugas"
}