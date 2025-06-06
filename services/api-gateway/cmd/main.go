package main

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	_ "github.com/DariaTarasek/diplom/services/api-gateway/docs"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/admin"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/auth"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/doctor"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/info"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/patient"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/statistics"
	"github.com/DariaTarasek/diplom/services/api-gateway/middleware"
	"github.com/DariaTarasek/diplom/services/api-gateway/perm"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title           Примастом API-Gateway
// @version         1.0
// @description     ПО для стоматологической клиники
// @host            localhost:8080
// @BasePath        /

func main() {
	r := gin.Default()
	authClient, err := clients.NewAuthClient("localhost:50052")
	if err != nil {
		log.Fatalf("Не удалось создать auth клиент: %s", err)
	}
	accessMiddleware := middleware.MakeAccessMiddleware(authClient)
	// Раздача статики (js/css/images)
	r.Static("/static", "./static/front")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Главная страница для пациента
	r.GET("/", func(c *gin.Context) {
		c.File("./static/templates/index.html")
	})

	// Страница авторизации персонала
	r.GET("/staff", func(c *gin.Context) {
		c.File("./static/templates/auth_doc.html")
	})

	adminClient, err := clients.NewAdminClient("localhost:50053")
	if err != nil {
		log.Fatalf("Не удалось создать admin клиент: %s", err)
	}

	storageClient, err := clients.NewStorageClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось создать storage клиент: %s", err)
	}

	patientClient, err := clients.NewPatientClient("localhost:50054")
	if err != nil {
		log.Fatalf("Не удалось создать patient клиент: %w", err)
	}

	doctorClient, err := clients.NewDoctorClient("localhost:50055")
	if err != nil {
		log.Fatalf("Не удалось создать doctor клиент: %s", err.Error())
	}

	statisticsClient, err := clients.NewStatisticsClient("localhost:50056")
	if err != nil {
		log.Fatalf("Не удалось создать statistics клиент: %s", err.Error())
	}

	htmlPages := []string{
		"index.html",
		"auth_doc.html",
		"registration.html",
		"auth.html",
		"employee_password_recovery.html",
		"password_recovery.html",
		"doctors.html",
		"appointment.html",
		"admins_schedule_management.html", // СТРАНИЦА АДМИНА! ЗДЕСЬ ДЛЯ ТЕСТОВ! ПОТОМ ПЕРЕНЕСТИ!
		"price_list.html",                 // СТРАНИЦА АДМИНА!
		"admins_admin_list.html",          // СТРАНИЦА АДМИНА!
		"admins_doctor_list.html",         // СТРАНИЦА АДМИНА!
		"admins_patient_list.html",        // СТРАНИЦА АДМИНА!
		"employee_registration.html",      // СТРАНИЦА АДМИНА!
		"registration_in_clinic.html",     // СТРАНИЦА АДМИНА!
		"administrator_account.html",      // СТРАНИЦА АДМИНА!
		"analytics.html",                  // СТРАНИЦА АДМИНА!

		//"doctor_account.html",       // СТРАНИЦА ВРАЧА!
		"doctors_consultation.html", // СТРАНИЦА ВРАЧА!

	}

	for _, page := range htmlPages {
		page := page // захват в замыкание
		r.GET("/"+page, func(c *gin.Context) {
			c.File("./static/templates/" + page)
		})
	}

	adminPages := []string{
		//"employee_registration.html",
		//"registration_in_clinic.html",
		//"admins_doctor_list.html",
		//"administrator_account.html",
		"administrator_profile.html",
	}
	for _, page := range adminPages {
		page := page // захват в замыкание
		r.GET("/"+page, accessMiddleware(perm.PermAdminPagesView), func(c *gin.Context) {
			c.File("./static/templates/" + page)
		})
	}

	doctorPages := []string{
		"doctor_account.html",
		"doctor_profile.html",
	}
	for _, page := range doctorPages {
		page := page // захват в замыкание
		r.GET("/"+page, accessMiddleware(perm.PermDoctorPagesView), func(c *gin.Context) {
			c.File("./static/templates/" + page)
		})
	}

	patientPages := []string{
		"patient_account.html",
		"patient_profile.html",
	}
	for _, page := range patientPages {
		page := page // захват в замыкание
		r.GET("/"+page, accessMiddleware(perm.PermPatientPagesView), func(c *gin.Context) {
			c.File("./static/templates/" + page)
		})
	}

	// REST API-группа
	api := r.Group("/api")

	registerHandler := auth.NewHandler(authClient)
	auth.RegisterRoutes(api, registerHandler)

	infoHandler := info.NewInfoHandler(storageClient)
	info.RegisterRoutes(api, infoHandler)

	adminHandler := admin.NewHandler(adminClient, authClient)
	admin.RegisterRoutes(api, adminHandler)

	patientHandler := patient.NewHandler(patientClient, accessMiddleware)
	patient.RegisterRoutes(api, patientHandler)

	doctorHandler := doctor.NewHandler(doctorClient, authClient)
	doctor.RegisterRoutes(api, doctorHandler)

	statisticsHandler := statistics.NewHandler(statisticsClient)
	statistics.RegisterRoutes(api, statisticsHandler)

	log.Println("Api-gateway запущен")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервис api-gateway: %v", err)
	}
}
