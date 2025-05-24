package main

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/info"
	"github.com/DariaTarasek/diplom/services/api-gateway/handlers/register"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	// Раздача статики (js/css/images)
	r.Static("/static", "./static/front")

	// Главная страница для пациента
	r.GET("/", func(c *gin.Context) {
		c.File("./static/templates/index.html")
	})

	// Страница авторизации персонала
	r.GET("/staff", func(c *gin.Context) {
		c.File("./static/templates/auth_doc.html")
	})

	authClient, err := clients.NewAuthClient("localhost:50052")
	if err != nil {
		log.Fatalf("Не удалось создать auth клиент: %s", err)
	}

	storageClient, err := clients.NewStorageClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось создать storage клиент: %s", err)
	}

	htmlPages := []string{"index.html", "auth_doc.html", "registration.html", "auth.html", "employee_registration.html"}

	for _, page := range htmlPages {
		page := page // захват в замыкание
		r.GET("/"+page, func(c *gin.Context) {
			c.File("./static/templates/" + page)
		})
	}

	// REST API-группа
	api := r.Group("/api")

	registerHandler := register.NewHandler(authClient)
	register.RegisterRoutes(api, registerHandler)

	infoHandler := info.NewInfoHandler(storageClient)
	info.RegisterRoutes(api, infoHandler)

	log.Println("Api-gateway запущен")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервис api-gateway: %v", err)
	}
}
