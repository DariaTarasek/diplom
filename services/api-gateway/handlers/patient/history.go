package patient

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	patientpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/patient"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getHistoryVisits godoc
// @Summary Получить историю посещений
// @Tags Пациент
// @Produce json
// @Success 200 {array} model.HistoryVisits
// @Failure 401 {object} map[string]string "Необходима авторизация"
// @Failure 403 {object} map[string]string "Недостаточно прав"
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /api/patient/history [get]
func (h *PatientHandler) getHistoryVisits(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен не найден"})
		return
	}
	resp, err := h.PatientClient.Client.GetHistoryVisits(c.Request.Context(), &patientpb.GetHistoryVisitsRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	historyVisits := make([]model.HistoryVisits, 0, len(resp.Visits))
	for _, item := range resp.Visits {
		historyVisit := model.HistoryVisits{
			ID:        int(item.Id),
			Date:      item.Date,
			DoctorID:  int(item.DoctorId),
			Doctor:    item.Doctor,
			Diagnose:  item.Diagnose,
			Treatment: item.Treatment,
		}
		historyVisits = append(historyVisits, historyVisit)
	}
	c.JSON(http.StatusOK, historyVisits)
}
