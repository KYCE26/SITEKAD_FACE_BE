package cuticontrollers

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
)

func CutiFileUpload(c *gin.Context) (string, error) {
	fileHeader, err := c.FormFile("suket")
	if err != nil {
		log.Printf("FormFile error: %v", err)
		if err == http.ErrMissingFile {
			return "", nil 
		}
		return "", fmt.Errorf("gagal membaca form file: %w", err)
	}
	log.Printf("File ditemukan: %s (%d bytes)", fileHeader.Filename, fileHeader.Size)

	// if err != nil {
	// 	if err == http.ErrMissingFile {
	// 		return "", nil 
	// 	}
	// 	return "", fmt.Errorf("gagal membaca form file: %w", err)
	// }
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

	uploadPath := os.Getenv("CUTI_PATH")
	if uploadPath == "" {
		uploadPath = "uploads/spl" 
	}

	destinationPath := filepath.Join(uploadPath, hashedFilename)
	if err := c.SaveUploadedFile(fileHeader, destinationPath); err != nil {
		return "", fmt.Errorf("gagal menyimpan file")
	}

	return hashedFilename, nil
}

func CreateCutiHandler(c *gin.Context) {

	userData, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	currentUser := userData.(models.Penempatan)

	hashedFilename, errFile := CutiFileUpload(c)
	if errFile != nil {
		log.Printf("Error during file upload: %v", errFile)
		c.JSON(http.StatusBadRequest, gin.H{"error": errFile.Error()})
		return
	}

	TanggalAwal := c.PostForm("tanggal_mulai")
	TanggalAkhir := c.PostForm("tanggal_selesai")
	alasan := c.PostForm("alasan")
	keterangan := c.PostForm("keterangan")

	if TanggalAwal == "" || TanggalAkhir == "" || alasan == "" {
		if hashedFilename != "" {
			uploadPath := os.Getenv("CUTI_PATH")
			if uploadPath == "" {
				uploadPath = "uploads/cuti"
			}
			os.Remove(filepath.Join(uploadPath, hashedFilename))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal mulai, tanggal selesai, dan alasan wajib diisi"})
		return
	}

	layout := "2006-01-02"
	tglMulai, _ := time.Parse(layout, TanggalAwal)
	tglSelesai, _ := time.Parse(layout, TanggalAkhir)
	if tglSelesai.Before(tglMulai) {
		if hashedFilename != "" {
			uploadPath := os.Getenv("CUTI_PATH")
			if uploadPath == "" {
				uploadPath = "uploads/cuti"
			}
			os.Remove(filepath.Join(uploadPath, hashedFilename))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal selesai tidak boleh sebelum tanggal mulai"})
		return
	}

	// now := time.Now()
	// tanggalHariIni := now.Format("2006-01-02")
	// jamSaatIni := now.Format("15:04:05")

	newCuti := models.Cuti{
		Penempatan_id: currentUser.Id,
		Tad_id:        currentUser.Pkwt.TadId,
		Cabang_id:     currentUser.Cabang_id,
		Lokasi_id:     currentUser.Lokasi_kerja_id,
		Jabatan_id:    currentUser.Jabatan_id,
		TglAwal:      TanggalAwal,
		TglAkhir:     TanggalAkhir,
		Alasan:        alasan,
		Keterangan:    keterangan,
		Suket:         hashedFilename,
	}

	if err := models.DB.Create(&newCuti).Error; err != nil {
		if hashedFilename != "" {
			uploadPath := os.Getenv("CUTI_PATH")
			if uploadPath == "" {
				uploadPath = "uploads/cuti"
			}
			os.Remove(filepath.Join(uploadPath, hashedFilename))
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengajuan cuti"})
		return
	}
	log.Printf("Hasil upload: filename=%s, err=%v", hashedFilename, errFile)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Pengajuan cuti berhasil diajukan!",
		"file_disimpan": hashedFilename,
	})
}