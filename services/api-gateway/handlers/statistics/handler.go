package statistics

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/DariaTarasek/diplom/services/api-gateway/perm"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	StatisticsClient *clients.StatisticsClient
	AccessMiddleware func(requiredPermission int32) gin.HandlerFunc
}

func NewHandler(statClient *clients.StatisticsClient, accessMiddleware func(requiredPermission int32) gin.HandlerFunc) *Handler {
	return &Handler{
		StatisticsClient: statClient,
		AccessMiddleware: accessMiddleware,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/statistics", h.AccessMiddleware(perm.PermSuperadminPagesView), h.GetAllStats)
}
