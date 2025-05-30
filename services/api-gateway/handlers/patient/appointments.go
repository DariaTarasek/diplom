package patient

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	patientpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/patient"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *PatientHandler) getAppointmentSlots(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("doctorId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slots, err := h.PatientClient.Client.GetAppointmentSlots(c.Request.Context(), &patientpb.GetAppointmentSlotsRequest{DoctorId: int32(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	scheduleSlots := make([]model.ScheduleEntry, 0, len(slots.Slots))
	for _, day := range slots.Slots {
		schSlot := model.ScheduleEntry{
			Label: day.Date,
			Slots: day.Slots,
		}
		scheduleSlots = append(scheduleSlots, schSlot)
	}

	c.JSON(http.StatusOK, scheduleSlots)
}

func (h *PatientHandler) addAppointment(c *gin.Context) {
	var req model.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dateStr := strings.Split(req.Date, "\n")
	date, err := time.Parse("02.01.2006", dateStr[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	appTime, err := time.Parse("15:04", req.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	birthDate, err := time.Parse("2006-01-02", req.PatientBirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment := &patientpb.Appointment{
		DoctorId:    int32(req.DoctorID),
		Date:        timestamppb.New(date),
		Time:        timestamppb.New(appTime),
		PatientId:   int32(derefUserID(req.PatientID)),
		SecondName:  req.PatientSecondName,
		FirstName:   req.PatientFirstName,
		Surname:     deref(req.PatientSurname),
		BirthDate:   timestamppb.New(birthDate),
		Gender:      req.PatientGender,
		PhoneNumber: req.PatientPhoneNumber,
	}
	_, err = h.PatientClient.Client.AddAppointment(c.Request.Context(), &patientpb.AddAppointmentRequest{Appointment: appointment})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Вы успешно записаны на прием!"})
}

func (h *PatientHandler) getUpcomingAppointments(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен не найден"})
		return
	}
	apps, err := h.PatientClient.Client.GetUpcomingAppointments(c.Request.Context(), &patientpb.GetUpcomingAppointmentsRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	upcoming := make([]model.UpcomingAppointment, 0, len(apps.Appointments))
	for _, app := range apps.Appointments {
		appointment := model.UpcomingAppointment{
			ID:        model.AppointmentID(app.Id),
			Date:      app.Date,
			Time:      app.Time,
			DoctorID:  model.UserID(app.DoctorId),
			Doctor:    app.Doctor,
			Specialty: app.Specialty,
		}
		upcoming = append(upcoming, appointment)
	}
	c.JSON(http.StatusOK, upcoming)
}

func (h *PatientHandler) UpdateAppointment(c *gin.Context) {
	var req model.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateDateStr := strings.Split(req.Date, "\n")
	updateDate, err := time.Parse("02.01.2006", updateDateStr[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateTime, err := time.Parse("15:04", req.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateApp := &patientpb.Appointment{
		Id:   int32(req.ID),
		Date: timestamppb.New(updateDate),
		Time: timestamppb.New(updateTime),
	}
	_, err = h.PatientClient.Client.UpdateAppointment(c.Request.Context(), &patientpb.UpdateAppointmentRequest{Appointment: updateApp})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись успешно обновлена"})
}

func (h *PatientHandler) CancelAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.PatientClient.Client.CancelAppointment(c.Request.Context(), &patientpb.GetByIDRequest{Id: int32(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись отменена"})
}

func derefUserID(u *model.UserID) model.UserID {
	if u == nil {
		return 0
	}
	return *u
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
