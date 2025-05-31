package doctor

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	doctorpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/doctor"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *DoctorHandler) GetTodayAppointments(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен не найден"})
		return
	}
	apps, err := h.DoctorClient.Client.GetTodayAppointments(c.Request.Context(), &doctorpb.GetTodayAppointmentsRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	todays := make([]model.TodayAppointment, 0, len(apps.Appointments))
	for _, app := range apps.Appointments {
		today := model.TodayAppointment{
			ID:        model.AppointmentID(app.Id),
			Date:      app.Date,
			Time:      app.Time,
			PatientID: model.UserID(app.PatientID),
			Patient:   app.Patient,
		}
		todays = append(todays, today)
	}

	c.JSON(http.StatusOK, todays)
}

func (h *DoctorHandler) GetUpcomingAppointments(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен не найден"})
		return
	}
	apps, err := h.DoctorClient.Client.GetUpcomingAppointments(c.Request.Context(), &doctorpb.GetUpcomingAppointmentsRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := convertRpcScheduleTable(apps.Schedule)

	c.JSON(http.StatusOK, res)
}

func convertRpcScheduleTable(rpcResp *doctorpb.ScheduleTable) model.ScheduleTable {
	result := model.ScheduleTable{
		Dates: rpcResp.Dates,
		Times: rpcResp.Times,
		Table: make(map[string]map[string]*model.UpcomingDoctorAppointment),
	}

	for _, cell := range rpcResp.Table {
		date := cell.Date
		time := cell.Time

		if result.Table[date] == nil {
			result.Table[date] = make(map[string]*model.UpcomingDoctorAppointment)
		}

		if cell.Appointment != nil {
			result.Table[date][time] = &model.UpcomingDoctorAppointment{
				ID:        model.AppointmentID(cell.Appointment.Id),
				PatientID: model.UserID(cell.Appointment.PatientId),
				Patient:   cell.Appointment.Patient,
			}
		} else {
			result.Table[date][time] = nil // Свободный слот
		}
	}

	return result
}
