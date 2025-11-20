package models

import "time"

type Penempatan struct {
	Id                int64     `gorm:"primaryKey" json:"id"`
	Nomor             string    `gorm:"type varchar(200)" json:"nomor"`
	Tgl               time.Time `gorm:"type date" json:"tgl"`
	Lokasi_kerja_id   int64     `gorm:"type int" json:"lokasi_kerja_id"`
	Pkwt_id           int64     `gorm:"type int" json:"pkwt_id"`
	Jabatan_id        int64     `gorm:"type int" json:"jabatan_id"`
	Created_at       time.Time `gorm:"type timestamp" json:"created_at"`
	Updated_at       time.Time `gorm:"type timestamp" json:"updated_at"`
	Created_by       int64     `gorm:"type int" json:"created_by"`
	Updated_by       int64     `gorm:"type int" json:"updated_by"`
	Is_deleted       int8      `gorm:"type tinyint" json:"is_deleted"`
	Cabang_id       int64     `gorm:"type int" json:"cabang_id"`
	Status_penempatan string    `gorm:"type varchar(100)" json:"status_penempatan"`
	Username         string    `gorm:"type varchar(255)" json:"username"`
	Password         string    `gorm:"type varchar(255)" json:"password"`
	Salt             string    `gorm:"type varchar(255)" json:"salt"`
	AndroidID       string    `gorm:"type varchar(100)" json:"android_id"`

	Pkwt	Pkwt      `gorm:"foreignKey:Pkwt_id"`
    Lokasi Lokasi `gorm:"foreignKey:Lokasi_kerja_id"`
    Cabang Cabang `gorm:"foreignKey:Cabang_id"`
}

func (Penempatan) TableName() string {
	return "penempatan"
}