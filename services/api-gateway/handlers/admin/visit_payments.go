package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetVisitPayments(c *gin.Context) {
	items, err := h.AdminClient.Client.GetUnconfirmedVisitPayments(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var visitPayments []model.VisitPayment
	for _, item := range items.VisitPayments {
		payment := model.VisitPayment{
			VisitID:   int(item.VisitId),
			Doctor:    item.Doctor,
			Patient:   item.Patient,
			CreatedAt: item.CreatedAt,
			Price:     int(item.Price),
		}
		visitPayments = append(visitPayments, payment)
	}
	c.JSON(http.StatusOK, visitPayments)
}
