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
