package models

type Cabang struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Kode    string `gorm:"type varchar(100)" json:"kode"`
	Nama    string `gorm:"type varchar(255)" json:"nama"`
	Zona	 string `gorm:"type varchar(10)" json:"zona"`
	Alamat  string `gorm:"type longtext" json:"alamat"`
	Telp   string `gorm:"type varchar(30)" json:"telp"`
	Is_deleted int8  `gorm:"type tinyint" json:"is_deleted"`
	Created_by int64  `gorm:"type bigint" json:"created_by"`
	Updated_by int64  `gorm:"type bigint" json:"updated_by"`
	Created_at string `gorm:"type:timestamp" json:"created_at"`
	Updated_at string `gorm:"type:timestamp" json:"updated_at"`
	Aktif   int8  `gorm:"type tinyint" json:"aktif"`
	Kepala_cabang string `gorm:"type varchar(200)" json:"kepala_cabang"`
	Kode_nitad string `gorm:"type varchar(5)" json:"kode_nitad"`
}

func (Cabang) TableName() string {
	return "cabang"
}	