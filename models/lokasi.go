package models

type Lokasi struct { 
	ID       uint   `gorm:"primaryKey" json:"id"`
	Kode    string `gorm:"type varchar(255)" json:"kode"`
	Nama    string `gorm:"type varchar(255)" json:"nama"`
	Alamat  string `gorm:"type varchar(255)" json:"alamat"`
	Telp   string `gorm:"type varchar(255)" json:"telp"`
	KotaId  int64  `gorm:"type bigint" json:"kota_id"`
	IsDeleted int8  `gorm:"type tinyint" json:"is_deleted"`
	CreatedAt string `gorm:"type:timestamp" json:"created_at"`
	CreatedBy int64  `gorm:"type bigint" json:"created_by"`
	UpdatedAt string `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy int64  `gorm:"type bigint" json:"updated_by"`
	MitraId  int64  `gorm:"type bigint" json:"mitra_id"`
	Zona	 string `gorm:"type varchar(255)" json:"zona"`
	Latitude  float64 `gorm:"type float" json:"latitude"`
	Longitude float64 `gorm:"type float" json:"longitude"`
	Kodeqr  string `gorm:"type varchar(255)" json:"kodeqr"`

}

func (Lokasi) TableName() string {
    return "lokasi_kerja"
}