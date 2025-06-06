package auth

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"time"
)

// @Summary Регистрация сотрудника
// @Tags Администратор
// @Accept json
// @Produce json
// @Param input body model.Employee true "Данные сотрудника"
// @Success 201 {object} gin.H
// @Failure 400,500 {object} gin.H
// @Router /api/employee-registe [post]
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

	var specsInt32 []int32
	for _, item := range employeeReq.Specs {
		specsInt32 = append(specsInt32, int32(item))
	}

	gRPCEmployee := &authpb.EmployeeData{
		FirstName:   employeeReq.FirstName,
		SecondName:  employeeReq.SecondName,
		Surname:     deref(employeeReq.Surname),
		PhoneNumber: deref(employeeReq.PhoneNumber),
		Email:       employeeReq.Email,
		Education:   education,
		Experience:  exp,
		Gender:      employeeReq.Gender,
		Role:        int32(employeeReq.Role),
		Specs:       specsInt32,
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

// @Summary Регистрация пациента
// @Tags Администратор
// @Accept json
// @Produce json
// @Param input body model.Patient true "Данные пациента"
// @Success 201 {object} gin.H
// @Failure 400,500 {object} gin.H
// @Router /api/register [post]
func (h *Handler) PatientRegister(c *gin.Context) {
	var patientReq model.Patient
	if err := c.ShouldBindJSON(&patientReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	gRPCUser := &authpb.UserData{
		Login:    patientReq.PhoneNumber,
		Password: patientReq.Password,
	}

	birthDate, err := time.Parse("2006-01-02", patientReq.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	gRPCPatient := &authpb.PatientData{
		FirstName:   patientReq.FirstName,
		SecondName:  patientReq.SecondName,
		Surname:     deref(patientReq.Surname),
		PhoneNumber: patientReq.PhoneNumber,
		Email:       deref(patientReq.Email),
		BirthDate:   timestamppb.New(birthDate),
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

// @Summary Регистрация пациента в клинике (без пароля)
// @Tags Администратор
// @Accept json
// @Produce json
// @Param input body model.PatientWithoutPassword true "Данные пациента"
// @Success 201 {object} gin.H
// @Failure 400,500 {object} gin.H
// @Router /api/register-in-clinic [post]
func (h *Handler) PatientRegisterInClinic(c *gin.Context) {
	var patientReq model.PatientWithoutPassword
	if err := c.ShouldBindJSON(&patientReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	gRPCUser := &authpb.UserData{
		Login:    deref(patientReq.PhoneNumber),
		Password: "",
	}

	birthDate, err := time.Parse("2006-01-02", patientReq.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	gRPCPatient := &authpb.PatientData{
		FirstName:   patientReq.FirstName,
		SecondName:  patientReq.SecondName,
		Surname:     deref(patientReq.Surname),
		PhoneNumber: deref(patientReq.PhoneNumber),
		Email:       deref(patientReq.Email),
		BirthDate:   timestamppb.New(birthDate),
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

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
