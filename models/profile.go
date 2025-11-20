package models

type Profile struct {
	Username   string `gorm:"type varchar(100)" json:"Username"`
	NomorInduk string `gorm:"type varchar(100)" json:"NITAD"`
	NamaLengkap string `gorm:"type varchar(200)" json:"Nama Lengkap"`
	NamaJabatan string `gorm:"type varchar(255)" json:"Jabatan"`
	NamaCabang string `gorm:"type varchar(255)" json:"Cabang"`
	NamaLokasi string `gorm:"type varchar(255)" json:"Lokasi"`
}