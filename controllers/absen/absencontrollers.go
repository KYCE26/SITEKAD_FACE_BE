package absen

import (
	"SITEKAD/helper"
	"SITEKAD/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllAbsen(c *gin.Context) {
	var absensi []models.Absensi

	if err := models.DB.Order("created_at desc").Find(&absensi).Limit(10).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Absen": absensi})
}

// Update Struct Payload agar fleksibel (QR atau Wajah)
type ScanPayload struct {
	Metode        string    `json:"metode"`         // "QR" atau "FACE"
	KodeQr        string    `json:"kodeqr"`         // Wajib jika QR
	FaceEmbedding []float64 `json:"face_embedding"` // Wajib jika FACE
	Latitude      float64   `json:"latitude" binding:"required"`
	Longitude     float64   `json:"longitude" binding:"required"`
	AndroidID     string    `json:"android_id" binding:"required"`
}

func ScanAbsensiHandler(c *gin.Context) {
	// 1. Bind JSON dari Android
	var payload ScanPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// 2. Ambil Data User Login
	userData, _ := c.Get("currentUser")
	currentUser := userData.(models.Penempatan)

	var targetLokasi models.LokasiPresensi
	var metodeAbsen string

	// =========================================================
	// LOGIKA SAKLAR: QR vs FACE
	// =========================================================

	if payload.Metode == "FACE" {
		// --- JALUR A: FACE RECOGNITION (MULTI ANGLE SUPPORT) ---
		metodeAbsen = "Face Recog"

		// A1. Validasi Input
		if len(payload.FaceEmbedding) != 192 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Data wajah rusak/tidak lengkap (Harus 192 dimensi)."})
			return
		}

		// A2. Ambil SEMUA Wajah Asli User dari Database
		var registeredFaces []models.UserFace
		if err := models.DB.Where("penempatan_id = ?", currentUser.Id).Find(&registeredFaces).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data wajah."})
			return
		}

		if len(registeredFaces) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wajah Anda belum didaftarkan. Silakan daftar wajah dulu."})
			return
		}

		// A3. Looping Mencari Kecocokan (Match)
		isMatch := false
		maxSimilarity := 0.0

		for _, face := range registeredFaces {
			var dbVector []float64
			// Skip jika data JSON di DB rusak
			if err := json.Unmarshal(face.Embedding, &dbVector); err != nil {
				continue
			}

			// Hitung Cosine Similarity
			score := helper.CosineSimilarity(payload.FaceEmbedding, dbVector)

			// Catat skor tertinggi untuk laporan error jika gagal
			if score > maxSimilarity {
				maxSimilarity = score
			}

			// Jika salah satu angle cocok > 80%, kita anggap Valid & Break loop
			if score >= 0.80 {
				isMatch = true
				break
			}
		}

		// A4. Putusan Akhir Wajah
		if !isMatch {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Wajah tidak dikenali! (Skor Tertinggi: %.1f%% - Butuh 80%%)", maxSimilarity*100),
			})
			return
		}

		// A5. Tentukan Target Lokasi (Default Lokasi User)
		err := models.DB.Where("id = ?", currentUser.Lokasi_kerja_id).First(&targetLokasi).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Lokasi kantor Anda tidak ditemukan di sistem."})
			return
		}

	} else {
		// --- JALUR B: QR CODE (Logika Lama) ---
		metodeAbsen = "QR Code"

		if payload.KodeQr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kode QR tidak boleh kosong."})
			return
		}

		// Validasi QR Code dan Lokasi
		err := models.DB.Where("kodeqr = ? AND penempatan_id = ?", payload.KodeQr, currentUser.Id).First(&targetLokasi).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "QR Code Salah atau Anda tidak terdaftar di lokasi ini!"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memvalidasi QR Code"})
			return
		}
	}

	// =========================================================
	// LOGIKA UMUM (RADIUS & SIMPAN DATA) - TIDAK BERUBAH
	// =========================================================

	// 1. Cek Radius (Geofencing)
	const radiusDiizinkanMeter = 300.0
	jarak := helper.Geolocation(payload.Latitude, payload.Longitude, targetLokasi.Latitude, targetLokasi.Longitude)

	if jarak > radiusDiizinkanMeter {
		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Anda berada di luar jangkauan (Jarak: %.0f meter)", jarak),
		})
		return
	}

	// 2. Persiapan Data Waktu
	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")
	jamSaatIni := now.Format("15:04:05")
	koordinatString := fmt.Sprintf("%f, %f", payload.Latitude, payload.Longitude)

	// 3. Cek Status Absen Hari Ini
	var absensi models.Absensi
	err := models.DB.Where("penempatan_id = ? AND tgl_absen = ?", currentUser.Id, tanggalHariIni).First(&absensi).Error

	switch err {
	case gorm.ErrRecordNotFound:
		// --- BELUM ABSEN -> CHECK-IN ---

		newAbsen := models.Absensi{
			Penempatan_id: currentUser.Id,
			Tad_id:        currentUser.Pkwt.TadId,
			Cabang_id:     currentUser.Cabang_id,
			Lokasi_id:     currentUser.Lokasi_kerja_id,
			Jabatan_id:    currentUser.Jabatan_id,
			Tgl_absen:     tanggalHariIni,
			Jam_masuk:     jamSaatIni,
			Kordmasuk:     koordinatString,
			Andid_masuk:   payload.AndroidID,
			Check:         tanggalHariIni + " " + jamSaatIni,
			Jenis:         &metodeAbsen, // Simpan jenis (QR/Face)
		}

		if err := models.DB.Create(&newAbsen).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi masuk"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Check-in BERHASIL (" + metodeAbsen + ") jam " + jamSaatIni})

	case nil:
		// --- SUDAH CHECK-IN -> CHECK-OUT ---

		if absensi.Jam_keluar != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah melakukan check-out hari ini"})
			return
		}

		// Validasi Durasi Kerja (Max 12 Jam)
		batasDurasi := 12 * time.Hour
		durasiSesi := time.Since(absensi.CreatedAt)
		if durasiSesi > batasDurasi {
			c.JSON(http.StatusForbidden, gin.H{"error": "Sesi kerja melebihi 12 jam. Hubungi admin."})
			return
		}

		// Handle Shift Malam
		var tanggalKeluar string
		hour := now.Hour()
		if hour >= 0 && hour < 6 {
			tanggalKeluar = now.AddDate(0, 0, -1).Format("2006-01-02")
		} else {
			tanggalKeluar = now.Format("2006-01-02")
		}

		// Update Checkout
		models.DB.Model(&absensi).Updates(models.Absensi{
			Tgl_keluar:   &tanggalKeluar,
			Jam_keluar:   &jamSaatIni,
			Kordkeluar:   &koordinatString,
			Andid_keluar: &payload.AndroidID,
			Check:        tanggalKeluar + " " + jamSaatIni,
		})
		c.JSON(http.StatusOK, gin.H{"message": "Check-out BERHASIL (" + metodeAbsen + ") jam " + jamSaatIni})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi masalah pada server: " + err.Error()})
	}
}

func GetHistoryUser(c *gin.Context) {
	userData, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	currentUser := userData.(models.Penempatan)
	var history []models.Absensi
	err := models.DB.Where("penempatan_id = ?", currentUser.Id).Order("tgl_absen DESC, jam_masuk DESC").Find(&history).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat absensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

func PrediksiCheckout(c *gin.Context) {
	userData, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	currentUser := userData.(models.Penempatan)
	var todayAbsen models.Absensi
	now := time.Now()
	tanggalHariIni := now.Format("2006-01-02")

	err := models.DB.Where("penempatan_id = ? AND tgl_absen = ?",
		currentUser.Id, tanggalHariIni).First(&todayAbsen).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Belum melakukan check-in hari ini",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data absensi",
		})
		return
	}
	if todayAbsen.Jam_keluar != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Sudah melakukan check-out hari ini",
			"jam_keluar": *todayAbsen.Jam_keluar,
		})
		return
	}
	history, err := helper.GetTrainingDataForUser(currentUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data historis",
		})
		return
	}
	if len(history) < 3 {
		c.JSON(http.StatusOK, gin.H{
			"message":               "Data historis tidak cukup untuk prediksi (minimal 3 data diperlukan)",
			"check_in":              todayAbsen.Jam_masuk,
			"prediction_available":  false,
			"historical_data_count": len(history),
		})
		return
	}
	predictedCheckout, err := helper.PredictCheckoutTime(history, todayAbsen.Jam_masuk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat prediksi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"check_in":          todayAbsen.Jam_masuk,
		"Prediksi Checkout": predictedCheckout,
		"Jumlah Absen":      len(history),
	})
}