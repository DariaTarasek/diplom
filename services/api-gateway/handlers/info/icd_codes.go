package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetICDCodes godoc
// @Summary Получить коды МКБ
// @Tags info
// @Produce json
// @Success 200 {array} model.ICDCode
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /api/icd-codes [get]
func (h *InfoHandler) GetICDCodes(c *gin.Context) {
	items, err := h.store.Client.GetICDCodes(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var codes []model.ICDCode
	for _, item := range items.IcdCode {
		code := model.ICDCode{
			ID:   int(item.Id),
			Code: item.Code,
			Name: item.Name,
		}
		codes = append(codes, code)
	}
	c.JSON(http.StatusOK, codes)
}
