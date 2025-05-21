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
		log.Fatalf("Ошибка подключения к БД: %s", err.Error())
	}
	defer conn.Close()

	migrationDSN := buildMigrationDSN()
	err = db.RunMigrations(migrationDSN)
	if err != nil {
		log.Fatalf("Миграция не была применена: %s", err.Error())
	}
	log.Println("Сервис БД запущен.")
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Переменная окружения %s не установлена.", key)
	}
	return value
}

func getDSNParams() map[string]string {
	result := make(map[string]string, 6)
	result["user"] = mustGetEnv("DB_USER")
	result["password"] = mustGetEnv("DB_PASSWORD")
	result["host"] = mustGetEnv("DB_HOST")
	result["port"] = mustGetEnv("DB_PORT")
	result["name"] = mustGetEnv("DB_NAME")
	result["sslMode"] = mustGetEnv("DB_SSLMODE")
	return result
}

// DSN - data source name
func buildDSN() string {
	params := getDSNParams()
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		params["user"], params["password"], params["host"], params["port"], params["name"], params["sslMode"])
}

func buildMigrationDSN() string {
	params := getDSNParams()
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		params["user"], params["password"], params["host"], params["port"], params["name"], params["sslMode"])
}
