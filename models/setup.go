package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. Coba load file .env (Khusus lokal)
	// Di Railway file ini ga ada, jadi pasti error.
	// Kita pakai '_ =' untuk MENGABAIKAN error-nya. Jangan di-Fatal!
	_ = godotenv.Load() 

	// 2. Ambil variabel dari System Environment (Railway Variable)
	dbURL := os.Getenv("DATABASE_URL")
	
	// 3. Baru cek di sini. Kalau variabelnya kosong, baru boleh panic.
	if dbURL == "" {
		log.Fatal("FATAL ERROR: Variable DATABASE_URL tidak ditemukan! Cek settings Railway.")
	}

	// 4. Konek ke Database
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal Terhubung ke Database: %v", err)
	}

	log.Println("Koneksi Database Berhasil.")
	DB = db
}