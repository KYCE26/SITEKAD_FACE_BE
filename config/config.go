package config

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Variable global untuk menyimpan key agar bisa diakses di controller/middleware
var JWT_KEY []byte

// Struct untuk data yang disimpan di dalam Token
type JWTClaims struct {
	Username string `json:"username"`
	// Kamu bisa tambah field lain di sini, misal: Role string
	jwt.RegisteredClaims
}

// Fungsi init berjalan otomatis saat aplikasi start
func init() {
	// 1. Coba load file .env (Khusus untuk Local Development di Laptop)
	// Di Railway, file ini biasanya tidak ada (masuk .gitignore), jadi akan error.
	// Kita abaikan error-nya dengan underscore (_), supaya aplikasi tidak mati.
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan. Menggunakan System Environment Variable (Mode Produksi/Railway).")
	}

	// 2. Ambil key dari Environment 
	// (Entah itu dari file .env lokal ATAU dari settingan 'Variables' di Railway)
	key := os.Getenv("JWT_KEY")
	
	// 3. Validasi Keamanan
	// Jika key tetap kosong (kelupaan setting), matikan aplikasi demi keamanan.
	if key == "" {
		log.Fatal("FATAL ERROR: JWT_KEY tidak ditemukan di environment variable! Pastikan sudah diset di .env atau Railway Variables.")
	}

	// 4. Simpan ke variable global sebagai byte slice
	JWT_KEY = []byte(key)
}