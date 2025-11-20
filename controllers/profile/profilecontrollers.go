package profilecontroller

import (

	"github.com/gin-gonic/gin"
	"net/http"
	"SITEKAD/models"
)


func GetUserProfile(c *gin.Context) {

	userData, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	
	currentUser := userData.(models.Penempatan)

	response := models.Profile{
		Username:     currentUser.Username,
		NomorInduk:   currentUser.Pkwt.Nitad,
		NamaLengkap:  currentUser.Pkwt.Tad.Nama,
		NamaJabatan:  currentUser.Pkwt.Jabatan.Nama, 
		NamaLokasi:   currentUser.Lokasi.Nama,
		NamaCabang:   currentUser.Cabang.Nama,
	}

	c.JSON(http.StatusOK, gin.H{"profile": response})
}