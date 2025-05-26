package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPassword(email string, password string, message string) error {
	// Настройки отправителя
	from := os.Getenv("EMAIL_ADDRESS")
	sourcePassword := os.Getenv("EMAIL_PASSWORD")

	// Настройки получателя
	to := []string{email}

	// SMTP-сервер
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Сообщение
	bytesMessage := []byte(message)

	// Аутентификация
	auth := smtp.PlainAuth("", from, sourcePassword, smtpHost)

	// Отправка письма
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, bytesMessage)
	if err != nil {
		fmt.Println("Ошибка при отправке:", err)
		return err
	}

	return nil
}
