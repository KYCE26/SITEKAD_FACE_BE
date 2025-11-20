package helper

import (
	"fmt"
	"strconv"
	"strings"

	"SITEKAD/models"

	"gonum.org/v1/gonum/stat"
)

func timeToMinutes(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0
	}
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	return float64(hours*60 + minutes)
}

func minutesToTime(minutes float64) string {
	hours := int(minutes / 60)
	mins := int(minutes) % 60
	return fmt.Sprintf("%02d:%02d", hours, mins)
}

func PredictCheckoutTime(history [][2]string, newCheckInTime string) (string, error) {
	if len(history) < 2 {
		return "", fmt.Errorf("minimal 2 data historis diperlukan untuk prediksi")
	}

	x := make([]float64, len(history))
	y := make([]float64, len(history))

	for i, record := range history {
		x[i] = timeToMinutes(record[0]) 
		y[i] = timeToMinutes(record[1]) 
	}

	alpha, beta := stat.LinearRegression(x, y, nil, false)
	newCheckInMinutes := timeToMinutes(newCheckInTime)
	predictedMinutes := alpha + beta*newCheckInMinutes

	return minutesToTime(predictedMinutes), nil
}

func GetTrainingDataForUser(userID int64) ([][2]string, error) {
	var recentAbsensi []struct {
		JamMasuk  string
		JamKeluar *string
	}

	err := models.DB.Table("absensi").
		Select("jam_masuk, jam_keluar").
		Where("penempatan_id = ? AND jam_keluar IS NOT NULL", userID).
		Order("id DESC").
		Limit(10).
		Scan(&recentAbsensi).Error

	if err != nil {
		return nil, err
	}

	historyData := make([][2]string, 0, len(recentAbsensi))
	for _, absen := range recentAbsensi {
		if absen.JamKeluar != nil {
			historyData = append(historyData, [2]string{absen.JamMasuk, *absen.JamKeluar})
		}
	}

	return historyData, nil
}
