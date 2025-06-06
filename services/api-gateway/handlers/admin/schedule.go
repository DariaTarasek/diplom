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

type updateClinicScheduleRequest struct {
	Schedule            []model.ClinicWeeklySchedule `json:"schedule"`
	SlotDurationMinutes int                          `json:"slot_duration_minutes"`
}

type updateDoctorScheduleRequest struct {
	Schedule            []model.DoctorWeeklySchedule `json:"schedule"`
	SlotDurationMinutes int                          `json:"slot_minutes"`
}

// UpdateClinicSchedule godoc
// @Summary Обновить недельное расписание клиники
// @Tags Администратор
// @Description Заменяет расписание клиники на новое, с заданной продолжительностью слота
// @Accept json
// @Produce json
// @Param schedule body updateClinicScheduleRequest true "Новое расписание клиники"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Некорректные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/clinic-schedule [put]
func (h *Handler) UpdateClinicSchedule(c *gin.Context) {
	var reqSchedule updateClinicScheduleRequest
	if err := c.ShouldBindJSON(&reqSchedule); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные: " + err.Error()})
		return
	}

	var schedule []*adminpb.WeeklyClinicSchedule
	for _, item := range reqSchedule.Schedule {
		if item.StartTime == "" {
			item.StartTime = "00:00"
		}
		start, err := time.Parse("15:04", item.StartTime)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось преобразовать время: " + err.Error()})
			return
		}
		if item.EndTime == "" {
			item.EndTime = "00:00"
		}
		end, err := time.Parse("15:04", item.EndTime)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось преобразовать время: " + err.Error()})
			return
		}
		day := &adminpb.WeeklyClinicSchedule{
			Id:                  int32(item.ID),
			Weekday:             int32(item.Weekday),
			StartTime:           timestamppb.New(start),
			EndTime:             timestamppb.New(end),
			SlotDurationMinutes: int32(reqSchedule.SlotDurationMinutes),
			IsDayOff:            item.IsDayOff,
		}
		schedule = append(schedule, day)
	}
	newSchedule := &adminpb.UpdateClinicWeeklyScheduleRequest{ClinicSchedule: schedule}
	_, err := h.AdminClient.Client.UpdateClinicWeeklySchedule(c.Request.Context(), newSchedule)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateDoctorSchedule godoc
// @Summary Обновить недельное расписание врача
// @Tags Администратор
// @Description Заменяет расписание врача на новое, с заданной продолжительностью слота
// @Accept json
// @Produce json
// @Param selectedDoctor path int true "ID врача"
// @Param schedule body updateDoctorScheduleRequest true "Новое расписание врача"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Некорректные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/doctor-schedule/{selectedDoctor} [put]
func (h *Handler) UpdateDoctorSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("selectedDoctor"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные: " + err.Error()})
		return
	}
	var reqSchedule updateDoctorScheduleRequest
	if err := c.ShouldBindJSON(&reqSchedule); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные: " + err.Error()})
		return
	}

	var schedule []*adminpb.WeeklyDoctorSchedule
	for _, item := range reqSchedule.Schedule {
		if item.StartTime == "" {
			item.StartTime = "00:00"
		}
		start, err := time.Parse("15:04", item.StartTime)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось преобразовать время: " + err.Error()})
			return
		}
		if item.EndTime == "" {
			item.EndTime = "00:00"
		}
		end, err := time.Parse("15:04", item.EndTime)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось преобразовать время: " + err.Error()})
			return
		}
		day := &adminpb.WeeklyDoctorSchedule{
			Id:                  int32(item.ID),
			DoctorId:            int32(id),
			Weekday:             int32(item.Weekday),
			StartTime:           timestamppb.New(start),
			EndTime:             timestamppb.New(end),
			SlotDurationMinutes: int32(reqSchedule.SlotDurationMinutes),
			IsDayOff:            item.IsDayOff,
		}
		schedule = append(schedule, day)
	}
	newSchedule := &adminpb.UpdateDoctorWeeklyScheduleRequest{DoctorSchedule: schedule}
	_, err = h.AdminClient.Client.UpdateDoctorWeeklySchedule(c.Request.Context(), newSchedule)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// AddClinicDailyOverride godoc
// @Summary Добавить переопределение расписания клиники на конкретный день
// @Tags Администратор
// @Description Устанавливает новое расписание или выходной день для клиники на указанную дату
// @Accept json
// @Produce json
// @Param override body model.ClinicDailyOverride true "Переопределение дня"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Некорректные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/clinic-overrides [post]
func (h *Handler) AddClinicDailyOverride(c *gin.Context) {
	var reqOverride model.ClinicDailyOverride
	if err := c.ShouldBindJSON(&reqOverride); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные: " + err.Error()})
		return
	}
	var isDayOff bool
	if reqOverride.IsDayOff == "off" {
		isDayOff = true
	} else {
		isDayOff = false
	}
	date, err := time.Parse("2006-01-02", reqOverride.Date)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if reqOverride.StartTime == "" {
		reqOverride.StartTime = "00:00"
	}
	if reqOverride.EndTime == "" {
		reqOverride.EndTime = "00:00"
	}
	start, err := time.Parse("15:04", reqOverride.StartTime)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	end, err := time.Parse("15:04", reqOverride.EndTime)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	_, err = h.AdminClient.Client.AddClinicDailyOverride(c.Request.Context(), &adminpb.AddClinicDailyOverrideRequest{
		Date:      timestamppb.New(date),
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
		IsDayOff:  isDayOff,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// AddDoctorDailyOverride godoc
// @Summary Добавить переопределение расписания врача на конкретный день
// @Tags Администратор
// @Description Устанавливает новое расписание или выходной день для врача на указанную дату
// @Accept json
// @Produce json
// @Param override body model.DoctorDailyOverride true "Переопределение дня врача"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Некорректные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/doctor-overrides [post]
func (h *Handler) AddDoctorDailyOverride(c *gin.Context) {
	var reqOverride model.DoctorDailyOverride
	if err := c.ShouldBindJSON(&reqOverride); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные: " + err.Error()})
		return
	}
	var isDayOff bool
	if reqOverride.IsDayOff == "off" {
		isDayOff = true
	} else {
		isDayOff = false
	}
	if reqOverride.StartTime == "" {
		reqOverride.StartTime = "00:00"
	}
	if reqOverride.EndTime == "" {
		reqOverride.EndTime = "00:00"
	}
	date, err := time.Parse("2006-01-02", reqOverride.Date)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	start, err := time.Parse("15:04", reqOverride.StartTime)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	end, err := time.Parse("15:04", reqOverride.EndTime)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	_, err = h.AdminClient.Client.AddDoctorDailyOverride(c.Request.Context(), &adminpb.AddDoctorDailyOverrideRequest{
		DoctorId:  int32(reqOverride.DoctorId),
		Date:      timestamppb.New(date),
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
		IsDayOff:  isDayOff,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

//func (h *Handler) GetUserRole(c *gin.Context) {
//	type Admin struct {
//		FirstName  string `json:"first_name"`
//		SecondName string `json:"second_name"`
//		Role       string `json:"role"`
//	}
//	MyAdmin := Admin{"Иван", "Иванов", "superadmin"}
//	c.JSON(http.StatusOK, MyAdmin)
//}
