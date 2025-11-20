package main

import (
	absencontroller "SITEKAD/controllers/absen"
	authcontroller "SITEKAD/controllers/auth"
	cuticontrollers "SITEKAD/controllers/cuti"
	facecontroller "SITEKAD/controllers/face"
	lemburcontrollers "SITEKAD/controllers/lembur"
	lokasicontroller "SITEKAD/controllers/lokasi"
	penugasancontrollers "SITEKAD/controllers/penugasan"
	profilecontroller "SITEKAD/controllers/profile"
	"SITEKAD/controllers/scheduler"
	"SITEKAD/middlewares"
	"SITEKAD/models"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func main() {
	models.ConnectDatabase()
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	//Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	v1 := router.Group("/v1")
	{
		v1.POST("/login", authcontroller.LoginHandler)
		v1.POST("/aktivasi", authcontroller.Aktivasi)
		v1.GET("/logout", authcontroller.Logout)

		api := v1.Group("/api")
		api.Use(middlewares.AuthMiddleware())
		{
			api.POST("/absensi", absencontroller.ScanAbsensiHandler)
			api.GET("/pcheck", absencontroller.PrediksiCheckout)
			api.GET("/histori", absencontroller.GetAllAbsen)
			api.GET("/lokasi", lokasicontroller.GetAllLokasi)
			api.GET("/uhistori", absencontroller.GetHistoryUser)
			api.GET("/profile", profilecontroller.GetUserProfile)
			api.POST("/lembur/start", lemburcontrollers.StartOvertimeHandler)
			api.PUT("/lembur/end", lemburcontrollers.EndOvertimeHandler)
			api.GET("/lembur/history", lemburcontrollers.GetHistoryLembur)
			api.POST("/patrol/start", penugasancontrollers.StartPatrolHandler)
			api.POST("/patrol/scan", penugasancontrollers.ScanCheckpointHandler)
			api.POST("/patrol/end", penugasancontrollers.EndPatrolHandler)
			api.POST("/cuti", cuticontrollers.CreateCutiHandler)
			api.POST("/scan", absencontroller.ScanAbsensiHandler)
			api.POST("/face/register", facecontroller.RegisterFaceHandler)
			api.GET("/face/status", facecontroller.CheckFaceStatusHandler)
		}
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(loc)
	s.Every(1).Hour().Do(scheduler.CleanupStalePatrols)
	s.StartAsync()
	log.Println("Automatic Scheduler Berhasil Dijalankan")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server Berhasil Berjalan Pada Port %s\n", port)

	router.Run(":" + port)
}
