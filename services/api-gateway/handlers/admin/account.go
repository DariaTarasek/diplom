package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetScheduleGrid godoc
// @Summary Получить расписание клиники
// @Tags Администратор
// @Description Возвращает таблицу расписания с приёмами, сгруппированную по дням и временным слотам
// @Produce json
// @Success 200 {object} model.AdminScheduleOverview
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/schedule-admin [get]
func (h *Handler) GetScheduleGrid(c *gin.Context) {
	gridResp, err := h.AdminClient.Client.GetClinicScheduleGrid(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем дни
	var scheduleDays []model.ScheduleDay
	for _, d := range gridResp.GetDays() {
		scheduleDays = append(scheduleDays, model.ScheduleDay{
			Date:    d.GetDate(),
			Weekday: d.GetWeekday(),
		})
	}

	// Преобразуем приёмы
	appointments := make(map[string]map[string][]model.AdminAppointment)

	for _, entry := range gridResp.GetAppointments() {
		date := entry.GetDate()
		time := entry.GetTime()

		if _, ok := appointments[date]; !ok {
			appointments[date] = make(map[string][]model.AdminAppointment)
		}

		appointments[date][time] = append(appointments[date][time], model.AdminAppointment{
			ID: int(entry.Id),
			Doctor: model.Person{
				ID:         model.UserID(entry.Doctor.Id),
				FirstName:  entry.Doctor.FirstName,
				SecondName: entry.Doctor.SecondName,
				Surname:    entry.Doctor.Surname,
				Specialty:  entry.Doctor.Specialty,
			},
			Patient: model.Person{
				ID:         model.UserID(entry.Patient.Id),
				FirstName:  entry.Patient.FirstName,
				SecondName: entry.Patient.SecondName,
				Surname:    entry.Patient.Surname,
				BirthDate:  entry.Patient.BirthDate,
				Gender:     entry.Patient.Gender,
				Phone:      entry.Patient.Phone,
			},
		})
	}

	// Отдаем JSON
	c.JSON(http.StatusOK, model.AdminScheduleOverview{
		Schedule: model.ScheduleMetadata{
			Days:      scheduleDays,
			TimeSlots: gridResp.GetTimeSlots(),
		},
		Appointments: appointments,
	})
}
