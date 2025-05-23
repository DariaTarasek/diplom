package main

import (
	"fmt"
	grpcserver "github.com/DariaTarasek/diplom/services/storage/grpc"
	"github.com/DariaTarasek/diplom/services/storage/internal/db"
	"github.com/DariaTarasek/diplom/services/storage/internal/store"
	pb "github.com/DariaTarasek/diplom/services/storage/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("не удалось начать слушать: %v", err)
	}

	s := grpc.NewServer()

	st := store.NewStore(conn) // твоя инициализация хранилища
	if err != nil {
		log.Fatalf("не удалось инициализировать store: %v", err)
	}

	server := &grpcserver.Server{
		Store: st,
	}

	pb.RegisterStorageServiceServer(s, server)

	log.Println("Storage gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
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
