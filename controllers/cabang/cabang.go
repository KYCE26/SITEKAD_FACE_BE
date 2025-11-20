package cabangcontroller

import (

	"SITEKAD/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Histori Absensi
func GetCabang(c *gin.Context) {
    var cabang []models.Cabang

    if err := models.DB.Order("created_at desc").Find(&cabang).Limit(10).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Cabang": cabang})
}
