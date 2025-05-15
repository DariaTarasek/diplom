package main

import (
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/db"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env файл не найден.")
	}

	dsn := buildDSN()

	conn, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %s", err)
	}
	defer conn.Close()

	log.Println("Сервис БД запущен.")
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Переменная окружения %s не установлена.", key)
	}
	return value
}

// DSN - data source name
func buildDSN() string {
	user := mustGetEnv("DB_USER")
	password := mustGetEnv("DB_PASSWORD")
	host := mustGetEnv("DB_HOST")
	port := mustGetEnv("DB_PORT")
	name := mustGetEnv("DB_NAME")
	sslMode := mustGetEnv("DB_SSLMODE")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslMode)
}
