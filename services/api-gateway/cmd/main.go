package main

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	handler "github.com/DariaTarasek/diplom/services/api-gateway/handlers"
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
		log.Fatalf("Не удалось создать auth клиент: %w", err)
	}

	htmlPages := []string{"index.html", "auth_doc.html", "registration.html", "auth.html"}

	for _, page := range htmlPages {
		page := page // захват в замыкание
		r.GET("/"+page, func(c *gin.Context) {
			c.File("./static/templates/" + page)
		})
	}

	// REST API-группа
	api := r.Group("/api")
	api.GET("/register", handler.Register)

	registerHandler := register.NewHandler(authClient)
	register.RegisterRoutes(api, registerHandler)

	log.Println("Api-gateway запущен")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервис api-gateway: %v", err)
	}
}
