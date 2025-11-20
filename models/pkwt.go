package models

type Pkwt struct {

	Id          int64  `gorm:"primaryKey" json:"id"`
	Nomor	   string `gorm:"type varchar(200)" json:"nomor"`
	Nitad   string `gorm:"type varchar(200)" json:"nitad"`
	NikMitra	string `gorm:"type varchar(100)" json:"nik_mitra"`
	MitraId	int64  `gorm:"type bigint" json:"mitra_id"`
	TadId  int64  `gorm:"type bigint" json:"tad_id"`
	TglAwal     string `gorm:"type date" json:"tgl_awal"`
	TglAkhir    string `gorm:"type date" json:"tgl_akhir"`
	TglBhk      string `gorm:"type date" json:"tgl_bhk"`
	KeteranganBhk string `gorm:"type varchar(200)" json:"keterangan_bhk"`
	JabatanId  int64  `gorm:"type bigint" json:"jabatan_id"`
	CabangId  int64  `gorm:"type bigint" json:"cabang_id"`
	CreatedAt   string `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   string `gorm:"type:timestamp" json:"updated_at"`
	CreatedBy   int64  `gorm:"type bigint" json:"created_by"`
	UpdatedBy   int64  `gorm:"type bigint" json:"updated_by"`
	IsDeleted   int8   `gorm:"type tinyint" json:"is_deleted"`
	StatusPkwt  string `gorm:"type varchar(200)" json:"status_pkwt"`
	KeteranganStatus string `gorm:"type varchar(255)" json:"keterangan_status"`
	JenisBhk string `gorm:"type varchar(255)" json:"jenis_bhk"`

	Jabatan Jabatan `gorm:"foreignKey:JabatanId"`
	Tad     Tad     `gorm:"foreignKey:TadId"`

}



func (Pkwt) TableName() string {
    return "pkwt"
}