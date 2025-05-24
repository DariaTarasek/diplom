package register

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"log"
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
	rg.POST("/employee-register", h.EmployeeRegister)
	//  сюда остальные
}

func (h *Handler) EmployeeRegister(c *gin.Context) {
	var employeeReq model.Employee
	if err := c.ShouldBindJSON(&employeeReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	password := "12345678cum"
	gRPCUser := &authpb.UserData{
		Login:    employeeReq.Email,
		Password: password,
	}

	var exp int32
	if employeeReq.Experience != nil {
		exp = int32(*employeeReq.Experience)
	}
	var education string
	if employeeReq.Education != nil {
		education = *employeeReq.Education
	}
	gRPCEmployee := &authpb.EmployeeData{
		FirstName:   employeeReq.FirstName,
		SecondName:  employeeReq.SecondName,
		Surname:     *employeeReq.Surname,
		PhoneNumber: *employeeReq.PhoneNumber,
		Email:       employeeReq.Email,
		Education:   education,
		Experience:  exp,
		Gender:      employeeReq.Gender,
		Role:        int32(employeeReq.Role),
	}

	gRPCEmployeeRequest := &authpb.EmployeeRegisterRequest{
		User:     gRPCUser,
		Employee: gRPCEmployee,
	}

	resp, err := h.AuthClient.Client.EmployeeRegister(c.Request.Context(), gRPCEmployeeRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": resp.UserId})
}
