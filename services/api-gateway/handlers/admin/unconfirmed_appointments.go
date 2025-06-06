package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"strconv"
	"time"
)

// @Summary Получить неподтверждённые записи
// @Tags Администратор
// @Description Возвращает список неподтверждённых записей на приём
// @Produce json
// @Success 200 {array} model.UnconfirmedAppointment
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/unconfirmed-appointments [get]
func (h *Handler) GetUnconfirmedAppointments(c *gin.Context) {
	items, err := h.AdminClient.Client.GetUnconfirmedAppointments(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var appointments []model.UnconfirmedAppointment
	for _, item := range items.Appointments {
		appt := model.UnconfirmedAppointment{
			ID:                int(item.Id),
			Doctor:            item.Doctor,
			PatientID:         int(item.PatientId),
			Date:              item.Date,
			Time:              item.Time,
			PatientFirstName:  item.FirstName,
			PatientSecondName: item.SecondName,
			PatientSurname:    item.Surname,
			PatientBirthDate:  item.BirthDate,
			Gender:            item.Gender,
			PhoneNumber:       item.PhoneNumber,
			Status:            item.Status,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
		}
		appointments = append(appointments, appt)
	}
	c.JSON(http.StatusOK, appointments)
}

// @Summary Обновить запись на приём
// @Tags Администратор
// @Description Обновляет дату, время и статус записи по ID
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Param appointment body model.UpdateAppointment true "Новая информация о записи"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/unconfirmed-appointments/{id} [put]
func (h *Handler) UpdateAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var appointment model.UpdateAppointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	date, err := time.Parse("02.01.2006", appointment.Date)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	apptTime, err := time.Parse("15:04", appointment.Time)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	updateAppt := &adminpb.UpdateAppointment{
		Id:        int32(id),
		Date:      timestamppb.New(date),
		Time:      timestamppb.New(apptTime),
		Status:    appointment.Status,
		UpdatedAt: timestamppb.New(appointment.UpdatedAt),
	}

	updateReq := &adminpb.UpdateAppointmentRequest{Appt: updateAppt}

	_, err = h.AdminClient.Client.UpdateAppointment(c.Request.Context(), updateReq)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
