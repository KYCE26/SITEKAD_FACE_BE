package authcontroller

import (
	"net/http"
	"time"

	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"SITEKAD/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type LoginPayload struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
    var payload LoginPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": "Input username dan password dibutuhkan"})
        return
    }

    var user models.Penempatan
    err := models.DB.Preload("Pkwt").Where("username = ?", payload.Username).First(&user).Error
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"Message": "Username atau Password Tidak Sesuai"})
        return
    }

 
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"Message": "Username atau Password Tidak Sesuai"})
        return
    }


    claims := jwt.MapClaims{
        "id":  user.Id, 
        "exp": time.Now().Add(time.Hour * 150).Unix(), 
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal membuat token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

type ActivationPayload struct {
	Nitad    string `json:"nitad" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	AndroidID string `json:"android_id" binding:"required"`
}

func Aktivasi(c *gin.Context) {
	var payload ActivationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Input tidak valid: " + err.Error()})
		return
	}

	var pkwt models.Pkwt
	err := models.DB.Where("nitad = ?", payload.Nitad).First(&pkwt).Error
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"Message": "Nitad tidak terdaftar"})
			return
		}
	
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Server error"})
		return
	}

	var existingUser models.Penempatan
	err = models.DB.Where("username = ?", payload.Username).First(&existingUser).Error

	if err == nil {
		
		c.JSON(http.StatusConflict, gin.H{"Message": "Username sudah digunakan, silakan pilih yang lain"})
		return
	} else if err != gorm.ErrRecordNotFound {
		
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal memvalidasi username"})
		return
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	result := models.DB.Model(&models.Penempatan{}).Where("pkwt_id = ?", pkwt.Id).Updates(models.Penempatan{
		Username: payload.Username,
		Password: string(hashPassword),
		AndroidID: payload.AndroidID,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengaktifkan akun"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Data pengguna untuk diaktifkan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Akun berhasil diaktifkan!"})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"Message": "Logout Berhasil!"})
}
