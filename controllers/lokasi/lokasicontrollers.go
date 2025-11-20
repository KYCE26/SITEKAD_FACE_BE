package lokasicontroller

import (

	"SITEKAD/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Histori Lokasi
func GetAllLokasi(c *gin.Context) {
    var lokasi []models.Lokasi

    if err := models.DB.Order("created_at desc").Find(&lokasi).Limit(10).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Lokasi": lokasi})
}