package register

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DoctorRegisterRequest struct {
	User   model.RegisterRequest `json:"user"`
	Doctor model.Doctor          `json:"doctor"`
}

type Handler struct {
	AuthClient *clients.AuthClient
}

func NewHandler(authClient *clients.AuthClient) *Handler {
	return &Handler{
		AuthClient: authClient,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/employee-register", h.DoctorRegister)
	// добавляй сюда остальные
}

func (h *Handler) DoctorRegister(c *gin.Context) {
	var req DoctorRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	gRPCUser := &authpb.UserData{
		Login:    req.User.Login,
		Password: req.User.Password,
	}

	exp := int32(*req.Doctor.Experience)
	gRPCDoctor := &authpb.DoctorData{
		FirstName:   req.Doctor.FirstName,
		SecondName:  req.Doctor.SecondName,
		Surname:     *req.Doctor.Surname,
		PhoneNumber: *req.Doctor.PhoneNumber,
		Email:       req.Doctor.Email,
		Education:   *req.Doctor.Education,
		Experience:  exp,
		Gender:      req.Doctor.Gender,
	}

	gRPCDocRequest := &authpb.DoctorRegisterRequest{
		User:   gRPCUser,
		Doctor: gRPCDoctor,
	}

	resp, err := h.AuthClient.Client.DoctorRegister(c.Request.Context(), gRPCDocRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": resp.UserId})
}
