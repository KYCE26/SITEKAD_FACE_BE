package models

import "time"

type Cuti struct {
	Id	int64	`gorm:"primaryKey" json:"id"`
	Penempatan_id int64   `gorm:"type int" json:"penempatan_id"`
	Tad_id        int64   `gorm:"type int" json:"tad_id"`
	Cabang_id     int64   `gorm:"type int" json:"cabang_id"`
	Lokasi_id     int64   `gorm:"type int" json:"lokasi_id"`
	Jabatan_id    int64   `gorm:"type int" json:"jabatan_id"`
	Suket         string	`gorm:"type varchar(255)" json:"suket"`
	TglAwal      string	`gorm:"type date" json:"tgl_awal"`
	TglAkhir     string	`gorm:"type date" json:"tgl_akhir"`
	Status	string	`gorm:"type varchar(50)" json:"status"`
	Alasan	string	`gorm:"type text" json:"alasan"`
	Keterangan	string	`gorm:"type text" json:"keterangan"`
	ApprovedBy	*int64	`gorm:"type int" json:"approved_by"`
	CreatedBy	*int64	`gorm:"type int" json:"created_by"`
	UpdatedBy	*int64	`gorm:"type int" json:"updated_by"`
	CreatedAt	time.Time	`gorm:"type timestamp" json:"created_at"`
	UpdatedAt	time.Time	`gorm:"type timestamp" json:"updated_at"`
	IsDeleted	int8	`gorm:"type tinyint" json:"is_deleted"`
}

func (Cuti) TableName() string {
	return "cuti"
}