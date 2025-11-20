package face

import (
	"SITEKAD/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct untuk validasi input dari Android
type RegisterFacePayload struct {
	Embedding []float64 `json:"embedding" binding:"required"`
}

func RegisterFaceHandler(c *gin.Context) {
	// 1. Ambil Data User yang sedang Login (Dari Middleware JWT)
	userData, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	currentUser := userData.(models.Penempatan)

	// 2. Validasi Input JSON
	var payload RegisterFacePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data wajah tidak valid: " + err.Error()})
		return
	}

	// 3. Validasi Dimensi Vektor (Wajib 192 untuk MobileFaceNet)
	if len(payload.Embedding) != 192 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dimensi vektor wajah salah (Harus 192)."})
		return
	}

	// 4. Konversi Array Float ke JSON String
	embeddingJSON, err := json.Marshal(payload.Embedding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data wajah"})
		return
	}

	// 5. SIMPAN DATA (APPEND MODE)
	// Kita selalu INSERT data baru, tidak menimpa yang lama.
	// Ini memungkinkan satu user memiliki banyak angle wajah.
	newFace := models.UserFace{
		PenempatanId: currentUser.Id,
		Name:         currentUser.Pkwt.Tad.Nama, // Nama diambil dari relasi TAD
		Embedding:    embeddingJSON,
	}

	if err := models.DB.Create(&newFace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data wajah"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Angle wajah berhasil disimpan!"})
}

// Fungsi Cek Status Wajah
func CheckFaceStatusHandler(c *gin.Context) {
	userData, _ := c.Get("currentUser")
	currentUser := userData.(models.Penempatan)

	var count int64
	// Hitung berapa banyak sampel wajah yang dimiliki user ini
	models.DB.Model(&models.UserFace{}).Where("penempatan_id = ?", currentUser.Id).Count(&count)

	// Dianggap terdaftar jika minimal punya 1 sampel
	c.JSON(http.StatusOK, gin.H{
		"is_registered": count > 0,
		"face_count":    count, // Info tambahan: jumlah angle yang tersimpan
	})
}