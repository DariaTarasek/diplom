package auth

import (
	"errors"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/DariaTarasek/diplom/services/api-gateway/sharederrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	requestCode struct {
		Phone string `json:"phone"`
	}
	verifyCode struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
)

func (h *Handler) requestCode(c *gin.Context) {
	var req requestCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.AuthClient.Client.RequestCode(c.Request.Context(), &authpb.GenerateCodeRequest{Phone: req.Phone})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Код отправлен на указанный номер"})
}

func (h *Handler) verifyCode(c *gin.Context) {
	var req verifyCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.AuthClient.Client.VerifyCode(c.Request.Context(), &authpb.VerifyCodeRequest{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		if errors.Is(err, sharederrors.ErrCodeInvalid) {
			c.JSON(http.StatusNotFound, gin.H{"error": "неверный введенный код"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "телефон подтвержден"})
}
