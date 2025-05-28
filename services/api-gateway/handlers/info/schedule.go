package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

func (h *InfoHandler) GetClinicWeeklySchedule(c *gin.Context) {
	items, err := h.store.Client.GetClinicWeeklySchedule(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var schedule []model.ClinicWeeklySchedule
	var slotMinutes int

	for _, item := range items.ClinicSchedule {
		day := model.ClinicWeeklySchedule{
			ID:                  int(item.Id),
			Weekday:             int(item.Weekday),
			StartTime:           FormatTime(item.StartTime.AsTime()),
			EndTime:             FormatTime(item.EndTime.AsTime()),
			SlotDurationMinutes: int(item.SlotDurationMinutes),
			IsDayOff:            !item.IsDayOff,
		}
		schedule = append(schedule, day)

		slotMinutes = int(item.SlotDurationMinutes)
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule":     schedule,
		"slot_minutes": slotMinutes,
	})
}

func (h *InfoHandler) GetDoctorWeeklySchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("selectedDoctor"))
	items, err := h.store.Client.GetDoctorWeeklySchedule(c.Request.Context(), &storagepb.GetScheduleByDoctorIdRequest{DoctorId: int32(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var schedule []model.DoctorWeeklySchedule
	var slotMinutes int

	for _, item := range items.DoctorSchedule {
		day := model.DoctorWeeklySchedule{
			ID:                  int(item.Id),
			DoctorID:            int(item.DoctorId),
			Weekday:             int(item.Weekday),
			StartTime:           FormatTime(item.StartTime.AsTime()),
			EndTime:             FormatTime(item.EndTime.AsTime()),
			SlotDurationMinutes: int(item.SlotDurationMinutes),
			IsDayOff:            !item.IsDayOff,
		}
		schedule = append(schedule, day)
		slotMinutes = int(item.SlotDurationMinutes)
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule":     schedule,
		"slot_minutes": slotMinutes,
	})
}

func (h *InfoHandler) GetClinicOverride(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Ожидается YYYY-MM-DD"})
		return
	}

	respOverride, err := h.store.Client.GetClinicOverride(c.Request.Context(),
		&storagepb.GetClinicOverrideRequest{Date: timestamppb.New(date)})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Переопределение не найдено"})
		return
	}

	override := model.ClinicDailyOverride{
		Date:     dateStr,
		IsDayOff: "work",
	}

	if respOverride.IsDayOff {
		override.IsDayOff = "off"
	}

	if respOverride.StartTime != nil {
		override.StartTime = respOverride.StartTime.AsTime().Format("15:04")
	}
	if respOverride.EndTime != nil {
		override.EndTime = respOverride.EndTime.AsTime().Format("15:04")
	}

	c.JSON(http.StatusOK, override)
}

func (h *InfoHandler) GetDoctorOverride(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Ожидается YYYY-MM-DD"})
		return
	}

	doctorIDStr := c.Param("doctor_id")
	id, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный идентификатор врача"})
		return
	}

	respOverride, err := h.store.Client.GetDoctorOverride(c.Request.Context(),
		&storagepb.GetDoctorOverrideRequest{
			DoctorId: int32(id),
			Date:     timestamppb.New(date),
		})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Переопределение для указанной даты и врача не найдено"})
		return
	}

	override := model.DoctorDailyOverride{
		DoctorId: int(respOverride.DoctorId),
		Date:     dateStr,
		IsDayOff: "work",
	}

	if respOverride.IsDayOff {
		override.IsDayOff = "off"
	}
	if respOverride.StartTime != nil {
		override.StartTime = respOverride.StartTime.AsTime().Format("15:04")
	}
	if respOverride.EndTime != nil {
		override.EndTime = respOverride.EndTime.AsTime().Format("15:04")
	}

	c.JSON(http.StatusOK, override)
}
