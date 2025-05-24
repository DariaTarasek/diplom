package register

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
)

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
	rg.POST("/register", h.PatientRegister)
	//  сюда остальные
}

func (h *Handler) EmployeeRegister(c *gin.Context) {
	var employeeReq model.Employee
	if err := c.ShouldBindJSON(&employeeReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	gRPCUser := &authpb.UserData{
		Login:    employeeReq.Email,
		Password: "",
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

func (h *Handler) PatientRegister(c *gin.Context) {
	var patientReq model.Patient
	if err := c.ShouldBindJSON(&patientReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	gRPCUser := &authpb.UserData{
		Login:    *patientReq.PhoneNumber,
		Password: patientReq.Password,
	}

	gRPCPatient := &authpb.PatientData{
		FirstName:   patientReq.FirstName,
		SecondName:  patientReq.SecondName,
		Surname:     *patientReq.Surname,
		PhoneNumber: *patientReq.PhoneNumber,
		Email:       *patientReq.Email,
		BirthDate:   timestamppb.New(patientReq.BirthDate),
		Gender:      patientReq.Gender,
	}

	gRPCPatientRequest := &authpb.PatientRegisterRequest{
		User:    gRPCUser,
		Patient: gRPCPatient,
	}

	resp, err := h.AuthClient.Client.PatientRegister(c.Request.Context(), gRPCPatientRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": resp.UserId})
}
