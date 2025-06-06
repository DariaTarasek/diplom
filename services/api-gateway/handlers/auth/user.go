package auth

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getPatient godoc
// @Summary Получить данные пациента
// @Description Возвращает информацию о текущем авторизованном пациенте на основе cookie access_token
// @Tags patient
// @Security ApiCookieAuth
// @Success 200 {object} model.PatientWithoutPassword
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/patient/me [get]
func (h *Handler) getPatient(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен не найден"})
		return
	}
	pbPatient, err := h.AuthClient.Client.GetPatient(c.Request.Context(), &authpb.GetPatientRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить пользователя"})
		return
	}
	patient := model.PatientWithoutPassword{
		ID:          int(pbPatient.Patient.UserId),
		FirstName:   pbPatient.Patient.FirstName,
		SecondName:  pbPatient.Patient.SecondName,
		Surname:     &pbPatient.Patient.Surname,
		PhoneNumber: &pbPatient.Patient.PhoneNumber,
		Email:       &pbPatient.Patient.Email,
		BirthDate:   pbPatient.Patient.BirthDate.AsTime().Format("2006-01-02"),
		Gender:      pbPatient.Patient.Gender,
	}
	c.JSON(http.StatusOK, patient)
}
