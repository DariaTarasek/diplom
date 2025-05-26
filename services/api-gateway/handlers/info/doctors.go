package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *InfoHandler) GetDoctors(c *gin.Context) {
	items, err := h.store.Client.GetDoctors(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
