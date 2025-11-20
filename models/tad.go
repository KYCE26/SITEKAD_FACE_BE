package models

type Tad struct {
	Id        int64  `gorm:"primaryKey" json:"id"`
	NoEktp	string `gorm:"type varchar(200)" json:"no_ektp"`
	Nama	string `gorm:"type varchar(200)" json:"nama"`
	JenisKelamin string `gorm:"type varchar(200)" json:"jenis_kelamin"`
	TempatLahir string `gorm:"type varchar(255)" json:"tempat_lahir"`
	TglLahir string `gorm:"type date" json:"tgl_lahir"`
	Alamat string `gorm:"type varchar(255)" json:"alamat"`
	Telp string `gorm:"type varchar(200)" json:"telp"`
	CabangId int64  `gorm:"type bigint" json:"cabang_id"`
	CreatedAt string `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp" json:"updated_at"`
	CreatedBy int64  `gorm:"type bigint" json:"created_by"`
	UpdatedBy int64  `gorm:"type bigint" json:"updated_by"`
	IsDeleted int8   `gorm:"type tinyint" json:"is_deleted"`
	Npwp string `gorm:"type varchar(200)" json:"npwp"`
	KodeJiwaNpwp string `gorm:"type varchar(2)" json:"kode_jiwa_npwp"`
}

func (Tad) TableName() string {
	return "tad"
}