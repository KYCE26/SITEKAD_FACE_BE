package penugasancontrollers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"SITEKAD/models"
	"gorm.io/gorm"
	"time"	
)

func StartPatrolHandler(c *gin.Context) {

	userData, _ := c.Get("currentUser")
	currentUser := userData.(models.Penempatan)

	var activePengerjaan models.PengerjaanTugas
	
	batasWaktu := time.Now().Add(-8 * time.Hour)

	err := models.DB.Where(
		"penempatan_id = ? AND status = ? AND waktu_mulai > ?",
		currentUser.Id,
		"berlangsung",
		batasWaktu,
	).First(&activePengerjaan).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":      "Anda sudah memiliki sesi patroli yang sedang berlangsung",
			"patroli_id": activePengerjaan.Ptid,
		})
		return
	}
	
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa sesi patroli"})
		return
	}

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Anda sudah memiliki sesi patroli yang sedang berlangsung",
			"ptid":  activePengerjaan.Ptid,
		})
		return
	}

	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa sesi patroli"})
		return
	}

	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")
	jamSaatIni := now.Format("15:04:05")

	newPengerjaan := models.PengerjaanTugas{
		PenempatanId: currentUser.Id,
		WaktuMulai:   tanggalHariIni + " " + jamSaatIni,
		Status:       "berlangsung",
	}

	if err := models.DB.Create(&newPengerjaan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai sesi patroli"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Sesi patroli berhasil dimulai",
		"ptid": newPengerjaan.Ptid,
	})
}

type ScanCheckpointPayload struct {
	PatroliID      int64  `json:"patroli_id" binding:"required"`
	CheckpointKode string `json:"checkpoint_kode" binding:"required"`
}

func ScanCheckpointHandler(c *gin.Context) {
	var payload ScanCheckpointPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	userData, _ := c.Get("currentUser")
	currentUser := userData.(models.Penempatan)

	var patroli models.PengerjaanTugas
	err := models.DB.Where(
		"ptid = ? AND penempatan_id = ? AND status = ?",
		payload.PatroliID, currentUser.Id, "berlangsung",
	).First(&patroli).Error

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sesi patroli tidak valid atau bukan milik Anda"})
		return
	}

	var checkpoint models.Checkpoint
	if err := models.DB.Where("kode_qr = ?", payload.CheckpointKode).First(&checkpoint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kode QR checkpoint tidak terdaftar"})
		return
	}

	var existingLog models.CekTugas
	err = models.DB.Where("ptid = ? AND cid = ?", patroli.Ptid, checkpoint.Cid).First(&existingLog).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Checkpoint ini sudah pernah Anda scan sebelumnya dalam sesi ini"})
		return
	}

	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")
	jamSaatIni := now.Format("15:04:05")

	newLog := models.CekTugas{
		Ptid:         patroli.Ptid,
		Cid:      checkpoint.Cid,
		WaktuScan:         tanggalHariIni + " " + jamSaatIni,
	}

	if err := models.DB.Create(&newLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat checkpoint"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Checkpoint '" + checkpoint.NamaLokasi + "' berhasil dicatat",
		"waktu_scan": newLog.WaktuScan,
	})
}

type EndPatrolPayload struct {
	PatroliID int64 `json:"patroli_id" binding:"required"`
}

func EndPatrolHandler(c *gin.Context) {
	var payload EndPatrolPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input patroli_id dibutuhkan"})
		return
	}

	userData, _ := c.Get("currentUser")
	currentUser := userData.(models.Penempatan)
	
	var patroli models.PengerjaanTugas
	err := models.DB.Where(
		"ptid = ? AND penempatan_id = ? AND status = ?",
		payload.PatroliID, currentUser.Id, "berlangsung",
	).First(&patroli).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada sesi patroli aktif yang cocok untuk diakhiri"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari sesi patroli"})
		return
	}

	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")
	jamSaatIni := now.Format("15:04:05")
	nowString := tanggalHariIni + " " + jamSaatIni
	result := models.DB.Model(&patroli).Updates(models.PengerjaanTugas{
		WaktuSelesai: &nowString,
		Status:       "selesai",
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan sesi patroli"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesi patroli berhasil diselesaikan"})
}

func GetPenugasan(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get Penugasan"})
}



