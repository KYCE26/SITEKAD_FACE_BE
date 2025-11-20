package models

import "time"

type Absensi struct {
	Id            int64   `gorm:"primaryKey" json:"id"`
	Penempatan_id int64   `gorm:"type int" json:"penempatan_id"`
	Tad_id        int64   `gorm:"type int" json:"tad_id"`
	Cabang_id     int64   `gorm:"type int" json:"cabang_id"`
	Lokasi_id     int64   `gorm:"type int" json:"lokasi_id"`
	Jabatan_id    int64   `gorm:"type int" json:"jabatan_id"`
	Tgl_absen     string  `gorm:"type date" json:"tgl_absen"`
	Jam_masuk     string  `gorm:"type time" json:"jam_masuk"`
	Tgl_keluar    *string `gorm:"type date" json:"tgl_keluar"`
	Jam_keluar    *string `gorm:"type time" json:"jam_keluar"`
	Check         string  `gorm:"type timestamp" json:"check_in"`
	Jenis         *string `gorm:"type varchar(255)" json:"jenis"`
	Kordmasuk     string  `gorm:"type varchar(255)" json:"kordmasuk"`
	Andid_masuk   string  `gorm:"type varchar(255)" json:"andid_masuk"`
	Kordkeluar    *string `gorm:"type varchar(255)" json:"kordkeluar"`
	Andid_keluar  *string `gorm:"type varchar(255)" json:"andid_keluar"`
	CreatedAt     time.Time `gorm:"type timestamp" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type timestamp" json:"updated_at"`
	Is_deleted    int8    `gorm:"type tinyint" json:"is_deleted"`
}

func (Absensi) TableName() string {
	return "absensi"
}