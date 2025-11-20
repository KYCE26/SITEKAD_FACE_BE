package models

type User struct {
	Id                int64  `gorm:"primaryKey" json:"id"`
	Username           string `gorm:"type varchar(255)" json:"username"`
	Email             string `gorm:"type varchar(255)" json:"email"`
	EmailVerifiedAt   string `gorm:"type timestamp" json:"email_verified_at"`
	Password          string `gorm:"type varchar(255)" json:"password"`
	RememberToken     string `gorm:"type varchar(100)" json:"remember_token"`
	CreatedAt         string `gorm:"type timestamp" json:"created_at"`
	UpdatedAt         string `gorm:"type timestamp" json:"updated_at"`
	RoleId           int64  `gorm:"type int" json:"role_id"`
	IsDeleted        int8   `gorm:"type tinyint" json:"is_deleted"`
	CreatedBy       int64  `gorm:"type int" json:"created_by"`
	UpdatedBy       int64  `gorm:"type int" json:"updated_by"`
	CabangId       int64  `gorm:"type int" json:"cabang_id"`
} 