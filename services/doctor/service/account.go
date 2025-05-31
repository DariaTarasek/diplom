package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/doctor/model"
	authpb "github.com/DariaTarasek/diplom/services/doctor/proto/auth"
	storagepb "github.com/DariaTarasek/diplom/services/doctor/proto/storage"
	"sort"
	"time"
	"unicode"
)

const (
	dayAhead   = 30
	weeksAhead = 4
	timeLayout = "15:04"
)

func (s *DoctorService) GetTodayAppointments(ctx context.Context, token string) ([]model.TodayAppointment, error) {
	userID, err := s.AuthClient.Client.GetUserID(ctx, &authpb.GetUserIDRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить userID: %w", err)
	}
	apps, err := s.StorageClient.Client.GetAppointmentsByDoctorID(ctx, &storagepb.GetAppointmentsByDoctorIDRequest{DoctorId: userID.UserId})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить записи: %w", err)
	}
	todays := make([]model.TodayAppointment, 0)
	now := time.Now().Truncate(24 * time.Hour)
	for _, app := range apps.Appointments {
		if app.Status == "cancelled" {
			continue
		}
		if !app.Date.AsTime().Equal(now) {
			continue
		}
		todayDate := app.Date.AsTime().Format("02.01.2006")
		todayTime := app.Time.AsTime().Format("15:04")
		patientName := fmt.Sprintf("%s %s.%s.", app.SecondName, getAndCapitalizeFirstLetter(app.FirstName), getAndCapitalizeFirstLetter(app.Surname))
		todayAppointment := model.TodayAppointment{
			ID:        model.AppointmentID(app.Id),
			Date:      todayDate,
			Time:      todayTime,
			PatientID: model.UserID(app.PatientId),
			Patient:   patientName,
		}
		todays = append(todays, todayAppointment)
	}
	return todays, nil
}

func (s *DoctorService) GetUpcomingAppointments(ctx context.Context, token string) (model.ScheduleTable, error) {
	doctorIDResp, err := s.AuthClient.Client.GetUserID(ctx, &authpb.GetUserIDRequest{Token: token})
	if err != nil {
		return model.ScheduleTable{}, fmt.Errorf("ошибка получения ID пользователя: %w", err)
	}

	scheduleResp, err := s.StorageClient.Client.GetDoctorWeeklySchedule(ctx, &storagepb.GetScheduleByDoctorIdRequest{
		DoctorId: doctorIDResp.UserId,
	})
	if err != nil {
		return model.ScheduleTable{}, fmt.Errorf("не удалось получить расписание врача: %w", err)
	}

	overridesResp, err := s.StorageClient.Client.GetDoctorOverrides(ctx, &storagepb.GetByIDRequest{Id: doctorIDResp.UserId})
	if err != nil {
		return model.ScheduleTable{}, fmt.Errorf("не удалось получить перегрузки врача: %w", err)
	}
	overridesMap := make(map[string]*storagepb.DoctorOverride)
	for _, override := range overridesResp.Override {
		if override.Date == nil {
			continue
		}
		dateStr := override.Date.AsTime().Format("02.01.2006")
		overridesMap[dateStr] = override
	}

	appointmentsResp, err := s.StorageClient.Client.GetAppointmentsByDoctorID(ctx, &storagepb.GetAppointmentsByDoctorIDRequest{
		DoctorId: doctorIDResp.UserId,
	})
	if err != nil {
		return model.ScheduleTable{}, fmt.Errorf("не удалось получить записи: %w", err)
	}

	// Сборка расписания врача
	var doctorSchedule []model.DoctorSchedule
	for _, item := range scheduleResp.DoctorSchedule {
		start, end := item.StartTime.AsTime(), item.EndTime.AsTime()
		slotDuration := int(item.SlotDurationMinutes)

		doctorSchedule = append(doctorSchedule, model.DoctorSchedule{
			ID:                  int(item.Id),
			DoctorID:            model.UserID(item.DoctorId),
			Weekday:             int(item.Weekday),
			StartTime:           &start,
			EndTime:             &end,
			SlotDurationMinutes: &slotDuration,
			IsDayOff:            &item.IsDayOff,
		})
	}

	// Преобразуем список приёмов
	appointmentsMap := map[string]map[string]*model.UpcomingAppointment{}
	for _, app := range appointmentsResp.Appointments {
		if app.Status == "cancelled" {
			continue
		}
		dateStr := fmt.Sprintf("%s\n(%s)", app.Date.AsTime().Format("02.01.2006"), weekdayToRus(app.Date.AsTime().Weekday()))
		timeStr := app.Time.AsTime().Format("15:04")

		if appointmentsMap[dateStr] == nil {
			appointmentsMap[dateStr] = make(map[string]*model.UpcomingAppointment)
		}

		patientID := model.UserID(app.PatientId)
		appointmentsMap[dateStr][timeStr] = &model.UpcomingAppointment{
			ID:        model.AppointmentID(app.Id),
			PatientID: patientID,
			Patient: fmt.Sprintf("%s %s.%s.",
				app.SecondName,
				getAndCapitalizeFirstLetter(app.FirstName),
				getAndCapitalizeFirstLetter(app.Surname)),
		}
	}

	// Подготовка таблицы
	result := model.ScheduleTable{
		Dates: []string{},
		Times: []string{},
		Table: map[string]map[string]*model.UpcomingAppointment{},
	}

	uniqueTimes := map[string]bool{}
	today := time.Now()
	weekdayToday := int(today.Weekday())
	if weekdayToday == 0 {
		weekdayToday = 7
	}
	monday := today.AddDate(0, 0, -weekdayToday+1)
	totalDays := weeksAhead * 7

	for i := 0; i < totalDays; i++ {
		date := monday.AddDate(0, 0, i)
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		dateStr := fmt.Sprintf("%s\n(%s)", date.Format("02.01.2006"), weekdayToRus(date.Weekday()))
		result.Dates = append(result.Dates, dateStr)
		result.Table[dateStr] = make(map[string]*model.UpcomingAppointment)

		override, hasOverride := overridesMap[date.Format("02.01.2006")]
		if hasOverride {
			if override.IsDayOff {
				// Перегрузка говорит, что выходной — пропускаем этот день
				continue
			}

			// Используем override-расписание
			start := time.Date(date.Year(), date.Month(), date.Day(),
				override.StartTime.AsTime().Hour(), override.StartTime.AsTime().Minute(), 0, 0, time.Local)
			end := time.Date(date.Year(), date.Month(), date.Day(),
				override.EndTime.AsTime().Hour(), override.EndTime.AsTime().Minute(), 0, 0, time.Local)

			for t := start; t.Before(end); t = t.Add(time.Duration(derefInt(doctorSchedule[1].SlotDurationMinutes)) * time.Minute) { // Можно заменить 30 на override-слот, если он появится
				timeStr := t.Format("15:04")
				uniqueTimes[timeStr] = true
				if appointment, exists := appointmentsMap[dateStr][timeStr]; exists {
					result.Table[dateStr][timeStr] = appointment
				} else {
					result.Table[dateStr][timeStr] = nil
				}
			}
			continue // Переход к следующему дню
		}

		for _, sched := range doctorSchedule {
			if sched.Weekday != weekday || *sched.IsDayOff {
				continue
			}

			start := time.Date(date.Year(), date.Month(), date.Day(), sched.StartTime.Hour(), sched.StartTime.Minute(), 0, 0, time.Local)
			end := time.Date(date.Year(), date.Month(), date.Day(), sched.EndTime.Hour(), sched.EndTime.Minute(), 0, 0, time.Local)

			for t := start; t.Before(end); t = t.Add(time.Duration(*sched.SlotDurationMinutes) * time.Minute) {
				timeStr := t.Format("15:04")
				uniqueTimes[timeStr] = true

				if appointment, exists := appointmentsMap[dateStr][timeStr]; exists {
					result.Table[dateStr][timeStr] = appointment
				} else {
					result.Table[dateStr][timeStr] = nil
				}
			}
		}
	}

	// Сортируем времена
	for timeStr := range uniqueTimes {
		result.Times = append(result.Times, timeStr)
	}
	sort.Strings(result.Times)

	return result, nil
}

func getAndCapitalizeFirstLetter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	return string(unicode.ToUpper(runes[0]))
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

func derefUserID(u *model.UserID) model.UserID {
	if u == nil {
		return 0
	}
	return *u
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
