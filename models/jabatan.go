package models

type Jabatan struct {
	Id   int64  `gorm:"primaryKey" json:"id"`
	JenisTadId int64  `gorm:"type int" json:"jenis_tad_id"`
	Nama string `gorm:"type varchar(255)" json:"nama"`
	Keterangan string `gorm:"type varchar(255)" json:"keterangan"`
	IsDeleted int8   `gorm:"type tinyint" json:"is_deleted"`
	CreatedAt string `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp" json:"updated_at"`
	CreatedBy int64  `gorm:"type bigint" json:"created_by"`
	UpdatedBy int64  `gorm:"type bigint" json:"updated_by"`
}

func (Jabatan) TableName() string {
	return "jabatan"
}