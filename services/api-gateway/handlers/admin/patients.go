package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetPatients godoc
// @Summary Получить список пациентов
// @Tags Администратор
// @Description Возвращает всех зарегистрированных пациентов
// @Produce json
// @Success 200 {array} model.PatientWithoutPassword
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/patients [get]
func (h *Handler) GetPatients(c *gin.Context) {
	items, err := h.AdminClient.Client.GetPatients(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var patients []model.PatientWithoutPassword
	for _, item := range items.Patients {
		patient := model.PatientWithoutPassword{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     &item.Surname,
			PhoneNumber: &item.PhoneNumber,
			Email:       &item.Email,
			Gender:      item.Gender,
			BirthDate:   item.BirthDate,
		}
		patients = append(patients, patient)
	}
	c.JSON(http.StatusOK, patients)
}

// UpdatePatient godoc
// @Summary Обновить данные пациента
// @Tags Администратор
// @Description Обновляет информацию о пациенте по ID
// @Accept json
// @Produce json
// @Param id path int true "ID пациента"
// @Param patient body model.Patient true "Новые данные пациента"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/patients/{id} [put]
func (h *Handler) UpdatePatient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var patientReq model.Patient
	if err := c.ShouldBindJSON(&patientReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdatePatientRequest := &adminpb.UpdatePatientRequest{
		UserId:      int32(id),
		FirstName:   patientReq.FirstName,
		Surname:     *patientReq.Surname,
		SecondName:  patientReq.SecondName,
		Email:       *patientReq.Email,
		BirthDate:   patientReq.BirthDate,
		PhoneNumber: patientReq.PhoneNumber,
		Gender:      patientReq.Gender,
	}

	_, err = h.AdminClient.Client.UpdatePatient(c.Request.Context(), UpdatePatientRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
