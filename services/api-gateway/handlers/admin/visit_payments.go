package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// @Summary Получить список неоплаченных посещений
// @Tags Администратор
// @Description Возвращает список неоплаченных посещений с материалами и услугами
// @Produce json
// @Success 200 {array} model.VisitPayment
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/completed-visits [get]
func (h *Handler) GetVisitPayments(c *gin.Context) {
	items, err := h.AdminClient.Client.GetUnconfirmedVisitPayments(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var visitPayments []model.VisitPayment
	for _, item := range items.VisitPayments {
		var materialsAndServices []model.MaterialsAndServices
		for _, mss := range item.MaterialsAndServices {
			ms := model.MaterialsAndServices{
				ID:       int(mss.Id),
				VisitID:  int(mss.VisitId),
				Item:     mss.Item,
				Quantity: int(mss.Quantity),
			}
			materialsAndServices = append(materialsAndServices, ms)
		}
		payment := model.VisitPayment{
			VisitID:              int(item.VisitId),
			Doctor:               item.Doctor,
			Patient:              item.Patient,
			CreatedAt:            item.CreatedAt,
			Price:                int(item.Price),
			MaterialsAndServices: materialsAndServices,
		}
		visitPayments = append(visitPayments, payment)
	}
	c.JSON(http.StatusOK, visitPayments)
}

// @Summary Обновить оплату за посещение
// @Tags Администратор
// @Description Обновляет информацию об оплате посещения по ID
// @Accept json
// @Produce json
// @Param id path int true "ID посещения"
// @Param payment body model.VisitPaymentUpdate true "Данные об оплате"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/completed-visits/{id} [put]
func (h *Handler) UpdateVisitPayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var payment model.VisitPaymentUpdate
	if err := c.ShouldBindJSON(&payment); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	updateVisitPayment := &adminpb.VisitPayment{
		VisitId: int32(id),
		Price:   int32(payment.Price),
		Status:  payment.Status,
	}

	updateReq := &adminpb.UpdateVisitPaymentRequest{Payment: updateVisitPayment}

	_, err = h.AdminClient.Client.UpdateVisitPayment(c.Request.Context(), updateReq)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
