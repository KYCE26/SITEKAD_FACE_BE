package scheduler

import (

	"SITEKAD/models" 
	"log"
	"time"
)


func CleanupStalePatrols() {
	log.Println("Memperbarui Data Tugas Patroli...")
	batasWaktu := time.Now().Add(-8 * time.Hour)
	var schedulePengerjaan []models.PengerjaanTugas
	err := models.DB.Where(
		"status = ? AND waktu_mulai < ?",
		"berlangsung",
		batasWaktu,
	).Find(&schedulePengerjaan).Error

    if err != nil {
        log.Printf("Gagal mencari sesi usang: %v\n", err)
        return
    }

	if len(schedulePengerjaan) == 0 {
		log.Println("Tidak ada sesi patroli usang yang ditemukan.")
		return
	}
	for _, pengerjaan := range schedulePengerjaan {
		log.Printf("Memperbarui sesi usang ID: %d", pengerjaan.Ptid)
		models.DB.Model(&pengerjaan).Update("status", "kedaluwarsa")
	}

	log.Printf("Selesai. %d sesi usang telah diperbarui.", len(schedulePengerjaan))
}