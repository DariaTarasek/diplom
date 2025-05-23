package handler

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	"net/http"

	// "github.com/DariaTarasek/diplom/services/api-gateway/proto"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	//grpcClient := client.GetClient()
	//resp, err := grpcClient.Register(c, &authpb.RegisterRequest{
	//	Username: req.Username,
	//	Password: req.Password,
	//	Email:    req.Email,
	//})
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC error: " + err.Error()})
	//	return
	//}

	c.JSON(http.StatusCreated, gin.H{"user_id": 1})
}
