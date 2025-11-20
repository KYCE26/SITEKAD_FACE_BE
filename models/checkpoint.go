package models

type Checkpoint struct {
	Cid	int64 `gorm:"primaryKey" json:"cid"`
	NamaLokasi	string `gorm:"type varchar(200)" json:"nama_lokasi"`
	KodeQr	string `gorm:"type varchar(200)" json:"kode_qr"`
}