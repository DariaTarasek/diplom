package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"time"
)

type updateClinicScheduleRequest struct {
	Schedule            []model.ClinicWeeklySchedule `json:"schedule"`
	SlotDurationMinutes int                          `json:"slot_duration_minutes"`
}

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

func (h *Handler) GetUserRole(c *gin.Context) {
	type Admin struct {
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		Role       string `json:"role"`
	}
	MyAdmin := Admin{"Ivan", "Ivanov", "superadmin"}
	c.JSON(http.StatusOK, MyAdmin)
}
