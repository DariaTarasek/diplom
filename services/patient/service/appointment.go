package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/patient/model"
	storagepb "github.com/DariaTarasek/diplom/services/patient/proto/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sort"
	"time"
)

const (
	dayAhead   = 30
	weeksAhead = 4
	timeLayout = "15:04"
)

func (s *PatientService) MakeDoctorAppointmentSlots(ctx context.Context, doctorID int) ([]model.ScheduleEntry, error) {
	resp, err := s.StorageClient.Client.GetDoctorWeeklySchedule(ctx, &storagepb.GetScheduleByDoctorIdRequest{DoctorId: int32(doctorID)})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить расписание врача: %w", err)
	}

	var doctorSchedule []model.DoctorSchedule
	for _, item := range resp.DoctorSchedule {
		start, end := item.StartTime.AsTime(), item.EndTime.AsTime()
		slotDuration := int(item.SlotDurationMinutes)
		doctorDay := model.DoctorSchedule{
			ID:                  int(item.Id),
			DoctorID:            model.UserID(item.DoctorId),
			Weekday:             int(item.Weekday),
			StartTime:           &start,
			EndTime:             &end,
			SlotDurationMinutes: &slotDuration,
			IsDayOff:            &item.IsDayOff,
		}
		doctorSchedule = append(doctorSchedule, doctorDay)
	}

	apps, err := s.StorageClient.Client.GetAppointmentsByDoctorID(ctx, &storagepb.GetAppointmentsByDoctorIDRequest{DoctorId: int32(doctorID)})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить существующие записи к врачу: %w", err)
	}

	var appointments []model.Appointment
	for _, app := range apps.Appointments {
		patientID := model.UserID(app.PatientId)
		appointment := model.Appointment{
			ID:                 model.AppointmentID(app.Id),
			DoctorID:           model.UserID(app.DoctorId),
			PatientID:          &patientID,
			Date:               app.Date.AsTime(),
			Time:               app.Time.AsTime(),
			PatientSecondName:  app.SecondName,
			PatientFirstName:   app.FirstName,
			PatientSurname:     &app.Surname,
			PatientBirthDate:   app.BirthDate.AsTime(),
			PatientGender:      app.Gender,
			PatientPhoneNumber: app.PhoneNumber,
			Status:             app.Status,
			CreatedAt:          app.CreatedAt.AsTime(),
			UpdatedAt:          app.UpdatedAt.AsTime(),
		}
		appointments = append(appointments, appointment)
	}

	busy := make(map[string]map[string]bool)
	for _, app := range appointments {
		dateStr := app.Date.Format("02.01.2006")
		timeStr := app.Time.Format("15:04")
		if busy[dateStr] == nil {
			busy[dateStr] = make(map[string]bool)
		}
		busy[dateStr][timeStr] = true
	}

	var tempSlots []struct {
		Date      time.Time
		DateLabel string
		Slots     []string
	}

	today := time.Now()
	weekdayToday := int(today.Weekday())
	if weekdayToday == 0 { // Sunday
		weekdayToday = 7
	}
	monday := today.AddDate(0, 0, -weekdayToday+1) // предыдущий (или сегодняшний) понедельник

	totalDays := weeksAhead * 7
	for i := 0; i < totalDays; i++ {
		date := monday.AddDate(0, 0, i)
		weekday := date.Weekday()

		slots := []string{}
		for _, day := range doctorSchedule {
			if time.Weekday(day.Weekday) == weekday && !*day.IsDayOff {
				startTime := time.Date(date.Year(), date.Month(), date.Day(), day.StartTime.Hour(), day.StartTime.Minute(), 0, 0, time.Local)
				endTime := time.Date(date.Year(), date.Month(), date.Day(), day.EndTime.Hour(), day.EndTime.Minute(), 0, 0, time.Local)

				for t := startTime; t.Before(endTime); t = t.Add(time.Duration(*day.SlotDurationMinutes) * time.Minute) {
					// Пропускаем слоты в прошлом, но день не исключаем
					if t.Before(time.Now()) {
						continue
					}

					timeStr := t.Format(timeLayout)
					busyKey := date.Format("02.01.2006")
					if busy[busyKey] != nil && busy[busyKey][timeStr] {
						continue
					}

					slots = append(slots, timeStr)
				}
			}
		}

		label := fmt.Sprintf("%s\n(%s)", date.Format("02.01.2006"), weekdayToRus(weekday))
		tempSlots = append(tempSlots, struct {
			Date      time.Time
			DateLabel string
			Slots     []string
		}{
			Date:      date,
			DateLabel: label,
			Slots:     slots,
		})
	}

	// Сортировка по дате
	sort.Slice(tempSlots, func(i, j int) bool {
		return tempSlots[i].Date.Before(tempSlots[j].Date)
	})

	// Финальная сборка отсортированного слайса
	var result []model.ScheduleEntry
	for _, slot := range tempSlots {
		result = append(result, model.ScheduleEntry{
			Label: slot.DateLabel,
			Slots: slot.Slots,
		})
	}

	log.Println(result)
	return result, nil

}

func (s *PatientService) AddAppointment(ctx context.Context, appointment model.Appointment) error {
	// TODO: Добавить сюда проверку, что GetAppointment(date, time) не существует
	appointmentPB := &storagepb.Appointment{
		DoctorId:    int32(appointment.DoctorID),
		Date:        timestamppb.New(appointment.Date),
		Time:        timestamppb.New(appointment.Time),
		PatientId:   int32(derefUserID(appointment.PatientID)),
		SecondName:  appointment.PatientSecondName,
		FirstName:   appointment.PatientFirstName,
		Surname:     deref(appointment.PatientSurname),
		BirthDate:   timestamppb.New(appointment.PatientBirthDate),
		Gender:      appointment.PatientGender,
		PhoneNumber: appointment.PatientPhoneNumber,
		Status:      "unconfirmed",
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
	}
	_, err := s.StorageClient.Client.AddAppointment(ctx, &storagepb.AddAppointmentRequest{Appointment: appointmentPB})
	if err != nil {
		return err
	}
	return nil
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

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefUserID(u *model.UserID) model.UserID {
	if u == nil {
		return 0
	}
	return *u
}
