package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPassword(email string, password string) error {
	// Настройки отправителя
	from := os.Getenv("EMAIL_ADDRESS")
	sourcePassword := os.Getenv("EMAIL_PASSWORD")

	// Настройки получателя
	to := []string{email}

	// SMTP-сервер (для Gmail)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Сообщение
	message := []byte(fmt.Sprintf("Subject: Регистрация в системе клиники!\r\n\r\nВы зарегистрированы в системе клиники!\nВаш пароль для входа: %s", password))

	// Аутентификация
	auth := smtp.PlainAuth("", from, sourcePassword, smtpHost)

	// Отправка письма
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Ошибка при отправке:", err)
		return err
	}

	return nil
}
