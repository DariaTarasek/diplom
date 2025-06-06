package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetDoctors godoc
// @Summary Получить всех врачей
// @Tags info
// @Produce json
// @Success 200 {array} model.Doctor
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /api/doctors [get]
func (h *InfoHandler) GetDoctors(c *gin.Context) {
	items, err := h.store.Client.GetDoctors(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var doctors []model.Doctor
	for _, item := range items.Doctors {
		exp := int(item.Experience)
		doctor := model.Doctor{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     &item.Surname,
			PhoneNumber: &item.PhoneNumber,
			Email:       item.Email,
			Education:   &item.Education,
			Experience:  &exp,
			Gender:      item.Gender,
		}
		doctors = append(doctors, doctor)
	}
	c.JSON(http.StatusOK, doctors)
}

// GetDoctorsBySpecialty godoc
// @Summary Получить врачей по специальности
// @Tags info
// @Produce json
// @Param specialty path int true "ID специальности"
// @Success 200 {array} model.Doctor
// @Failure 400 {object} map[string]string "Некорректный ID специальности"
// @Failure 500 {object} map[string]string "Ошибка получения данных"
// @Router /api/doctors/{specialty} [get]
func (h *InfoHandler) GetDoctorsBySpecialty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("specialty"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос " + err.Error()})
		return
	}
	items, err := h.store.Client.GetDoctorsBySpecID(c.Request.Context(), &storagepb.GetDoctorBySpecIDRequest{SpecId: int32(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список врачей по специальности " + err.Error()})
		return
	}
	var doctors []model.Doctor
	for _, item := range items.Doctors {
		exp := int(item.Experience)
		doctor := model.Doctor{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     &item.Surname,
			PhoneNumber: &item.PhoneNumber,
			Email:       item.Email,
			Education:   &item.Education,
			Experience:  &exp,
			Gender:      item.Gender,
		}
		doctors = append(doctors, doctor)
	}
	c.JSON(http.StatusOK, doctors)
}
