package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"log"
	"time"
	"unicode"
)

const weeksAhead = 4

func (s *AdminService) GetClinicScheduleGrid(ctx context.Context) (model.AdminScheduleOverview, error) {
	scheduleResp, err := s.StorageClient.Client.GetClinicWeeklySchedule(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AdminScheduleOverview{}, fmt.Errorf("не удалось получить расписание клиники: %w", err)
	}

	overrideResp, err := s.StorageClient.Client.GetClinicOverrides(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AdminScheduleOverview{}, fmt.Errorf("не удалось получить перегрузки клиники: %w", err)
	}
	overridesMap := make(map[string]*storagepb.ClinicOverride)
	for _, override := range overrideResp.Overrides {
		if override.Date == nil {
			continue
		}
		dateStr := override.Date.AsTime().Format("02.01.2006")
		overridesMap[dateStr] = override
	}

	appointmentsResp, err := s.StorageClient.Client.GetAppointments(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AdminScheduleOverview{}, fmt.Errorf("не удалось получить записи клиники: %w", err)
	}

	appointmentsMap := map[string]map[string][]*storagepb.Appointment{}
	doctorIDSet := map[int]struct{}{}
	for _, app := range appointmentsResp.Appointments {
		if app.Status == "cancelled" {
			continue
		}
		dateStr := app.Date.AsTime().Format("02.01.2006")
		timeStr := app.Time.AsTime().Format("15:04")

		if appointmentsMap[dateStr] == nil {
			appointmentsMap[dateStr] = make(map[string][]*storagepb.Appointment)
		}
		appointmentsMap[dateStr][timeStr] = append(appointmentsMap[dateStr][timeStr], app)

		if app.DoctorId != 0 {
			doctorIDSet[int(app.DoctorId)] = struct{}{}
		}
	}

	doctorInfoMap := map[int]model.Person{}
	for doctorID := range doctorIDSet {
		doctorResp, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: int32(doctorID)})
		if err != nil {
			log.Printf("не удалось получить врача по ID %d: %v", doctorID, err)
			continue
		}
		specsResp, err := s.StorageClient.Client.GetDoctorSpecsByDoctorId(ctx, &storagepb.GetByIdRequest{Id: int32(doctorID)})
		if err != nil {
			log.Printf("не удалось получить спеки врача по ID %d: %v", doctorID, err)
			continue
		}
		doctorInfoMap[doctorID] = model.Person{
			ID:         model.UserID(doctorResp.Doctor.UserId),
			FirstName:  doctorResp.Doctor.FirstName,
			SecondName: doctorResp.Doctor.SecondName,
			Surname:    doctorResp.Doctor.Surname,
			Specialty:  fetchSpecsIntoString(specsResp.Specs),
		}
	}

	result := model.AdminScheduleOverview{
		Schedule: model.ScheduleMetadata{
			Days:      []model.ScheduleDay{},
			TimeSlots: []string{},
		},
		Appointments: map[string]map[string][]model.AdminAppointment{},
	}

	// Получаем slotDuration из первого расписания (все одинаковые)
	var defaultSlotDuration int32 = 30
	if len(scheduleResp.ClinicSchedule) > 0 {
		if scheduleResp.ClinicSchedule[0].SlotDurationMinutes > 0 {
			defaultSlotDuration = scheduleResp.ClinicSchedule[0].SlotDurationMinutes
		}
	}

	today := time.Now()
	weekdayToday := int(today.Weekday())
	if weekdayToday == 0 {
		weekdayToday = 7
	}
	monday := today.AddDate(0, 0, -weekdayToday+1)
	totalDays := weeksAhead * 7

	clinicSchedule := map[int]*storagepb.WeeklyClinicSchedule{}
	for _, sched := range scheduleResp.ClinicSchedule {
		clinicSchedule[int(sched.Weekday)] = sched
	}

	var globalStart, globalEnd *time.Time
	updateGlobalTimeRange := func(start, end time.Time) {
		if globalStart == nil || start.Before(*globalStart) {
			globalStart = &start
		}
		if globalEnd == nil || end.After(*globalEnd) {
			globalEnd = &end
		}
	}

	for i := 0; i < totalDays; i++ {
		date := monday.AddDate(0, 0, i)
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		dateStr := date.Format("02.01.2006")
		weekdayStr := weekdayToRus(date.Weekday())
		result.Schedule.Days = append(result.Schedule.Days, model.ScheduleDay{
			Date:    dateStr,
			Weekday: weekdayStr,
		})
		result.Appointments[dateStr] = map[string][]model.AdminAppointment{}

		if override, ok := overridesMap[dateStr]; ok {
			if override.IsDayOff {
				continue
			}
			start := override.StartTime.AsTime()
			end := override.EndTime.AsTime()
			updateGlobalTimeRange(start, end)

			for t := start; t.Before(end); t = t.Add(time.Duration(defaultSlotDuration) * time.Minute) {
				timeStr := t.Format("15:04")
				entries := []model.AdminAppointment{}
				if apps, ok := appointmentsMap[dateStr][timeStr]; ok {
					for _, app := range apps {
						doctor, _ := doctorInfoMap[int(app.DoctorId)]
						entries = append(entries, model.AdminAppointment{
							ID:     int(app.Id),
							Doctor: doctor,
							Patient: model.Person{
								ID:         model.UserID(app.PatientId),
								SecondName: app.SecondName,
								FirstName:  app.FirstName,
								Surname:    app.Surname,
								BirthDate:  app.BirthDate.AsTime().Format("02.01.2006"),
								Gender:     app.Gender,
								Phone:      app.PhoneNumber,
							},
						})
					}
				}
				result.Appointments[dateStr][timeStr] = entries
			}
			continue
		}

		sched, ok := clinicSchedule[weekday]
		if !ok || sched.IsDayOff {
			continue
		}
		start := sched.StartTime.AsTime()
		end := sched.EndTime.AsTime()
		updateGlobalTimeRange(start, end)

		for t := start; t.Before(end); t = t.Add(time.Duration(defaultSlotDuration) * time.Minute) {
			timeStr := t.Format("15:04")
			entries := []model.AdminAppointment{}
			if apps, ok := appointmentsMap[dateStr][timeStr]; ok {
				for _, app := range apps {
					doctor, _ := doctorInfoMap[int(app.DoctorId)]
					entries = append(entries, model.AdminAppointment{
						ID:     int(app.Id),
						Doctor: doctor,
						Patient: model.Person{
							ID:         model.UserID(app.PatientId),
							SecondName: app.SecondName,
							FirstName:  app.FirstName,
							Surname:    app.Surname,
							BirthDate:  app.BirthDate.AsTime().Format("02.01.2006"),
							Gender:     app.Gender,
							Phone:      app.PhoneNumber,
						},
					})
				}
			}
			result.Appointments[dateStr][timeStr] = entries
		}
	}

	// Генерация всех возможных слотов
	if globalStart == nil || globalEnd == nil {
		start, _ := time.Parse("15:04", "08:00")
		end, _ := time.Parse("15:04", "20:00")
		globalStart = &start
		globalEnd = &end
	}
	for t := *globalStart; !t.After(*globalEnd); t = t.Add(time.Duration(defaultSlotDuration) * time.Minute) {
		result.Schedule.TimeSlots = append(result.Schedule.TimeSlots, t.Format("15:04"))
	}

	return result, nil
}

func weekdayToRus(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "Пн"
	case time.Tuesday:
		return "Вт"
	case time.Wednesday:
		return "Ср"
	case time.Thursday:
		return "Чт"
	case time.Friday:
		return "Пт"
	case time.Saturday:
		return "Сб"
	case time.Sunday:
		return "Вс"
	}
	return ""
}

func getAndCapitalizeFirstLetter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	return string(unicode.ToUpper(runes[0]))
}

func fetchSpecsIntoString(specs []int32) string {
	var specsStr string
	for _, spec := range specs {
		specsStr += fetchSpec(spec) + ", "
	}
	return specsStr[:len(specsStr)-2]
}

func fetchSpec(spec int32) string {
	switch spec {
	case model.Therapist:
		return "терапевт"
	case model.Surgeon:
		return "хирург"
	case model.Orthopedist:
		return "ортопед"
	default:
		return ""

	}
}
