package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func TestUpdateAppointmentIntegration(t *testing.T) {
	// URL вашего работающего сервера. Убедитесь, что сервер запущен.
	appointmentID := 17

	// Подготовка данных
	payload := map[string]interface{}{
		"id":         appointmentID,
		"date":       "04.06.2025", // формат "дд.мм.гггг"
		"time":       "14:30",      // формат "чч:мм"
		"status":     "confirmed",
		"updated_at": time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	// Формируем PUT-запрос
	req, err := http.NewRequest(http.MethodPut, baseURL+"/api/unconfirmed-appointments/17", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверка статуса ответа
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEmployeeRegister_Success(t *testing.T) {
	payload := map[string]interface{}{
		"firstName":       "Иван",
		"secondName":      "Иванов",
		"surname":         "Иванович",
		"phone":           "79999999999",
		"email":           "ivan123@example.com",
		"education":       "SPbSUT",
		"experience":      3,
		"gender":          "м",
		"role":            3,
		"specializations": []int{1, 2, 3},
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/api/employee-register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось отправить запрос: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("Ожидается код 201, получен код %d", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидается код 201, получен код %d", resp.StatusCode)
	}
}

func TestEmployeeRegister_BadRequest(t *testing.T) {
	// Отправляем некорректный JSON
	invalidJSON := `{"firstName": "Иван", "email": 123}` // email должен быть строкой

	resp, err := http.Post(baseURL+"/api/employee-register", "application/json", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Не удалось отправить запрос: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("Ожидается код 400, получен код %d", resp.StatusCode)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Ожидается код 400, получен код %d", resp.StatusCode)
	}
}
