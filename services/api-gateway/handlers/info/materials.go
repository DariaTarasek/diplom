package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *InfoHandler) GetMaterials(c *gin.Context) {
	items, err := h.store.Client.GetMaterials(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var materials []model.Material
	for _, item := range items.Materials {
		material := model.Material{
			ID:    int(item.Id),
			Name:  item.Name,
			Price: int(item.Price),
		}
		materials = append(materials, material)
	}
	c.JSON(http.StatusOK, gin.H{"materials": materials})
}
