package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
