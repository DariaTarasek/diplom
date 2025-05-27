package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
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
