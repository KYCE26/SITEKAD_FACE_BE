package models

type LokasiPresensi struct {
	
	ID       int64   `gorm:"primaryKey" json:"id"`
	CabangId  int64  `gorm:"type bigint" json:"cabang_id"`
	NamaLokasiPresensi    string `gorm:"type varchar(255)" json:"nama_lokasi_presensi"`
	LokasiKerjaId int64  `gorm:"type bigint" json:"lokasi_kerja_id"`
	PenempatanId	int64  `gorm:"type bigint" json:"penempatan_id"`
	Pkwt_id           int64     `gorm:"type int" json:"pkwt_id"`
	TadId  int64  `gorm:"type bigint" json:"tad_id"`
	Jabatan_id        int64     `gorm:"type int" json:"jabatan_id"`
	Longitude float64 `gorm:"type float" json:"longitude"`
	Latitude  float64 `gorm:"type float" json:"latitude"`
	Zona	 string `gorm:"type varchar(255)" json:"zona"`
	Kodeqr  string `gorm:"type varchar(255)" json:"kodeqr"`
	Fileqr  string `gorm:"type varchar(255)" json:"fileqr"`
	IsDeleted int8  `gorm:"type tinyint" json:"is_deleted"`
	CreatedAt string `gorm:"type:timestamp" json:"created_at"`
	CreatedBy int64  `gorm:"type bigint" json:"created_by"`
	UpdatedAt string `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy int64  `gorm:"type bigint" json:"updated_by"`	

}

func (LokasiPresensi) TableName() string {
    return "lokasi_kerja"
}