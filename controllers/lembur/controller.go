package lemburcontrollers

import (
	"SITEKAD/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func handleFileUpload(c *gin.Context) (string, error) {
	fileHeader, err := c.FormFile("spl_file")
	if err != nil {
		return "", fmt.Errorf("file SPL wajib diunggah")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file")
	}
	defer file.Close()
	ext := filepath.Ext(fileHeader.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Format File Tidak Sesuai! Hanya file jpg, png, pdf Yang Diizinkan!"})
		return "", fmt.Errorf("format file tidak diizinkan")
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file untuk validasi")
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("gagal mereset file reader")
	}


	mimeType := http.DetectContentType(buffer)
	allowedMimes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}

	if !allowedMimes[mimeType] {
		return "", fmt.Errorf("tipe konten file terdeteksi tidak valid: %s", mimeType)
	}

	extension := filepath.Ext(fileHeader.Filename)
	stringToHash := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
	hasher := sha256.New()
	hasher.Write([]byte(stringToHash))
	hashedFilename := hex.EncodeToString(hasher.Sum(nil)) + extension
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "uploads/spl"
	}

	destinationPath := filepath.Join(uploadPath, hashedFilename)
	if err := c.SaveUploadedFile(fileHeader, destinationPath); err != nil {
		return "", fmt.Errorf("gagal menyimpan file")
	}

	return hashedFilename, nil
}

func StartOvertimeHandler(c *gin.Context) {
	userData, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	currentUser := userData.(models.Penempatan)
	var existingLembur models.Lembur
	twelveHoursAgo := time.Now().Add(-12 * time.Hour)
	err := models.DB.Where("penempatan_id = ? AND jam_keluar IS NULL AND created_at > ?", currentUser.Id, twelveHoursAgo).First(&existingLembur).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah memiliki sesi lembur yang aktif."})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi sesi lembur"})
		return
	}

	hashedFilename, errFile := handleFileUpload(c)
	if errFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errFile.Error()})
		return
	}

	latitude := c.PostForm("latitude")
	longitude := c.PostForm("longitude")
	androidID := c.PostForm("android_id")
	if len(longitude) > 50 || len(androidID) > 50 || len(latitude) > 50 || len(longitude) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input Tidak Valid!"})
		return
	}

	if latitude == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon Aktifkan Izin Lokasi!"})
		return
	}
	if longitude == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon Aktifkan Izin Lokasi!"})
		return
	}
	if androidID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon Berikan Izin Akses Perangkat Agar Dapat Melakukan Absensi!"})
        return
    }
	koordinat := fmt.Sprintf("%s, %s", latitude, longitude)
	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")
	jamSaatIni := now.Format("15:04:05")

	newLembur := models.Lembur{
		Penempatan_id: currentUser.Id,
		Tad_id:        currentUser.Pkwt.TadId,
		Cabang_id:     currentUser.Cabang_id,
		Lokasi_id:     currentUser.Lokasi_kerja_id,
		Jabatan_id:    currentUser.Jabatan_id,
		Spl:           hashedFilename,
		Tgl_absen:     tanggalHariIni,
		Jam_masuk:     jamSaatIni,
		Kordmasuk:     koordinat,
		Andid_masuk:   androidID,
		Check:         tanggalHariIni + " " + jamSaatIni,
	}

	if errDb := models.DB.Create(&newLembur).Error; errDb != nil {
		
		uploadPath := os.Getenv("UPLOAD_PATH")
		if uploadPath == "" {
			uploadPath = "./uploads"
		}
		os.Remove(filepath.Join(uploadPath, hashedFilename)) 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data lembur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Sesi lembur berhasil dimulai",
		"file_disimpan": hashedFilename,
	})
}

type EndOvertimePayload struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	AndroidID string  `json:"android_id" binding:"required"`
}

func EndOvertimeHandler(c *gin.Context) {

	userData, _ := c.Get("currentUser")
	currentUser := userData.(models.Penempatan)

	var payload EndOvertimePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return

	}

	koordinatString := fmt.Sprintf("%f, %f", payload.Latitude, payload.Longitude)
	var lembur models.Lembur
	err := models.DB.Where("penempatan_id = ? AND jam_keluar IS NULL", currentUser.Id).Last(&lembur).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada sesi lembur aktif yang ditemukan untuk diakhiri"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari data lembur"})
		return
	}

	if payload.AndroidID != lembur.Andid_masuk {
		log.Printf("ID MISMATCH -> Payload: '[%s]' vs DB: '[%s]'", payload.AndroidID, lembur.Andid_masuk)
		c.JSON(http.StatusForbidden, gin.H{"error": "Perangkat yang digunakan untuk check-out berbeda dengan saat check-in"})
		return
	}

	if err := models.DB.Model(&lembur).Where("jam_keluar IS NULL").First(&lembur).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Sesi lembur sudah diakhiri sebelumnya"})
		return
	}
	

	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")
	jamSaatIni := now.Format("15:04:05")
	result := models.DB.Model(&lembur).Updates(models.Lembur{
		Tgl_keluar:   &tanggalHariIni,
		Jam_keluar:   &jamSaatIni,
		Kordkeluar:   &koordinatString,
		Andid_keluar: &payload.AndroidID,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan sesi lembur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesi lembur berhasil diakhiri pada jam " + jamSaatIni})
}

func GetHistoryLembur(c *gin.Context) {
    userData, exists := c.Get("currentUser")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
        return
    }
    currentUser := userData.(models.Penempatan)
    var history []models.Lembur
    err := models.DB.Where("penempatan_id = ?", currentUser.Id).Order("tgl_absen DESC, jam_masuk DESC").Find(&history).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat lembur"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"history": history})
}
