package statistics

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	statspb "github.com/DariaTarasek/diplom/services/api-gateway/proto/statistics"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllStats godoc
// @Summary      Получить общую статистику по клинике
// @Description  Возвращает агрегированную информацию: общее число пациентов, визитов, топ-услуги, среднюю загруженность и выручку по врачам, распределение по возрастным группам и др.
// @Tags         Статистика
// @Produce      json
// @Success      200  {object}  model.AllStats
// @Failure      500  {object}  map[string]string
// @Router       /api/statistics [get]
func (h *Handler) GetAllStats(c *gin.Context) {
	items, err := h.StatisticsClient.Client.GetAllStats(c.Request.Context(), &statspb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Преобразование TopServices
	var topServices []model.TopService
	for _, item := range items.TopServices {
		topServices = append(topServices, model.TopService{
			Name:  item.Name,
			Count: int(item.Count),
		})
	}

	// Преобразование DoctorAvgVisit
	var doctorAvgVisit []model.DoctorAvgVisit
	for _, item := range items.DocAvgVisits {
		doctorAvgVisit = append(doctorAvgVisit, model.DoctorAvgVisit{
			Doctor:          item.Doctor,
			AvgWeeklyVisits: item.AvgWeeklyVisits,
		})
	}

	// Преобразование DoctorCheckStat
	var doctorCheckStat []model.DoctorCheckStat
	for _, item := range items.DoctorCheck {
		doctorCheckStat = append(doctorCheckStat, model.DoctorCheckStat{
			Doctor:   item.Doctor,
			AvgCheck: item.AvgCheck,
		})
	}

	// Преобразование DoctorUniquePatients
	var doctorUniquePatients []model.DoctorUniquePatients
	for _, item := range items.DoctorPatients {
		doctorUniquePatients = append(doctorUniquePatients, model.DoctorUniquePatients{
			DoctorID:       item.Doctor,
			UniquePatients: int(item.UniquePatients),
		})
	}

	// Преобразование AgeGroupStat
	var ageGroupStat []model.AgeGroupStat
	for _, item := range items.AgeGroups {
		ageGroupStat = append(ageGroupStat, model.AgeGroupStat{
			AgeGroup: item.AgeGroup,
			Percent:  item.Percent,
		})
	}

	stats := model.AllStats{
		TotalPatients:        int(items.TotalPatients),
		TotalVisits:          int(items.TotalVisits),
		TopServices:          topServices,
		DoctorAvgVisit:       doctorAvgVisit,
		DoctorCheckStat:      doctorCheckStat,
		DoctorUniquePatients: doctorUniquePatients,
		AgeGroupStat:         ageGroupStat,
		NewPatientsThisMonth: int(items.NewPatientsThisMonth),
		AvgVisitPerPatient:   items.AvgVisitPerPatient,
		TotalIncome:          items.TotalIncome,
		MonthlyIncome:        items.MonthlyIncome,
		ClinicAvgCheck:       items.ClinicAvgCheck,
	}

	c.JSON(http.StatusOK, stats)
}
