package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/patient/model"
	authpb "github.com/DariaTarasek/diplom/services/patient/proto/auth"
	storagepb "github.com/DariaTarasek/diplom/services/patient/proto/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sort"
	"time"
	"unicode"
)

const (
	dayAhead   = 30
	weeksAhead = 4
	timeLayout = "15:04"
)

type Override struct {
	StartTime time.Time
	EndTime   time.Time
	SlotMins  int
	IsDayOff  bool
}

func (s *PatientService) MakeDoctorAppointmentSlots(ctx context.Context, doctorID int) ([]model.ScheduleEntry, error) {
	resp, err := s.StorageClient.Client.GetDoctorWeeklySchedule(ctx, &storagepb.GetScheduleByDoctorIdRequest{DoctorId: int32(doctorID)})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить расписание врача: %w", err)
	}
	docOverrides, err := s.StorageClient.Client.GetDoctorOverrides(ctx, &storagepb.GetByIDRequest{Id: int32(doctorID)})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить переопределения врача: %w", err)
	}
	clinicOverrides, err := s.StorageClient.Client.GetClinicOverrides(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить переопределения клиники: %w", err)
	}

	clinicOverridesMap := make(map[string]Override)
	for _, o := range clinicOverrides.Overrides {
		dateStr := o.Date.AsTime().Format("02.01.2006")
		clinicOverridesMap[dateStr] = Override{
			StartTime: o.StartTime.AsTime(),
			EndTime:   o.EndTime.AsTime(),
			SlotMins:  int(o.SlotDurationMinutes),
			IsDayOff:  o.IsDayOff,
		}
	}

	doctorOverridesMap := make(map[string]Override)
	for _, o := range docOverrides.Override {
		dateStr := o.Date.AsTime().Format("02.01.2006")
		doctorOverridesMap[dateStr] = Override{
			StartTime: o.StartTime.AsTime(),
			EndTime:   o.EndTime.AsTime(),
			SlotMins:  int(o.SlotDurationMinutes),
			IsDayOff:  o.IsDayOff,
		}
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
		if app.Status == "cancelled" {
			continue
		}
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
	monday := today.AddDate(0, 0, -int(today.Weekday())+1)

	totalDays := weeksAhead * 7
	for i := 0; i < totalDays; i++ {
		date := monday.AddDate(0, 0, i)
		dateStr := date.Format("02.01.2006")
		dateBusyStr := date.Format("02.01.2006")
		weekday := date.Weekday()

		var slots []string

		clinicOverride, hasClinic := clinicOverridesMap[dateStr]
		doctorOverride, hasDoctor := doctorOverridesMap[dateStr]

		switch {
		// Приоритет: выходной клиники => нет слотов
		case hasClinic && clinicOverride.IsDayOff:
			slots = append(slots, []string{}...)

		// Если врач взял выходной — игнорируем, даже если клиника работает
		case hasDoctor && doctorOverride.IsDayOff:
			slots = append(slots, []string{}...)

		// Если есть обе перегрузки
		case hasClinic && hasDoctor:
			clinicStart := time.Date(date.Year(), date.Month(), date.Day(), clinicOverride.StartTime.Hour(), clinicOverride.StartTime.Minute(), 0, 0, time.Local)
			clinicEnd := time.Date(date.Year(), date.Month(), date.Day(), clinicOverride.EndTime.Hour(), clinicOverride.EndTime.Minute(), 0, 0, time.Local)

			doctorStart := time.Date(date.Year(), date.Month(), date.Day(), doctorOverride.StartTime.Hour(), doctorOverride.StartTime.Minute(), 0, 0, time.Local)
			doctorEnd := time.Date(date.Year(), date.Month(), date.Day(), doctorOverride.EndTime.Hour(), doctorOverride.EndTime.Minute(), 0, 0, time.Local)

			if !doctorStart.Before(clinicStart) && !doctorEnd.After(clinicEnd) {
				// Врач внутри клиники — используем врача
				slots = append(slots, generateSlots(doctorStart, doctorEnd, doctorOverride.SlotMins, busy[dateBusyStr])...)
			} else {
				// Врач вне рамок — используем клинику
				slots = append(slots, generateSlots(clinicStart, clinicEnd, clinicOverride.SlotMins, busy[dateBusyStr])...)
			}

		// Только перегрузка врача
		case hasDoctor:
			start := time.Date(date.Year(), date.Month(), date.Day(), doctorOverride.StartTime.Hour(), doctorOverride.StartTime.Minute(), 0, 0, time.Local)
			end := time.Date(date.Year(), date.Month(), date.Day(), doctorOverride.EndTime.Hour(), doctorOverride.EndTime.Minute(), 0, 0, time.Local)
			slots = append(slots, generateSlots(start, end, doctorOverride.SlotMins, busy[dateBusyStr])...)

		// Только перегрузка клиники
		case hasClinic:
			start := time.Date(date.Year(), date.Month(), date.Day(), clinicOverride.StartTime.Hour(), clinicOverride.StartTime.Minute(), 0, 0, time.Local)
			end := time.Date(date.Year(), date.Month(), date.Day(), clinicOverride.EndTime.Hour(), clinicOverride.EndTime.Minute(), 0, 0, time.Local)
			slots = append(slots, generateSlots(start, end, clinicOverride.SlotMins, busy[dateBusyStr])...)

		// Никаких перегрузок — обычное расписание
		default:
			log.Printf("weekday=%d, doctorSchedule count=%d", weekday, len(doctorSchedule))
			for _, day := range doctorSchedule {
				log.Printf("checking day: weekday=%d, target=%d, isDayOff=%v", day.Weekday, weekday, derefBool(day.IsDayOff))
				log.Printf("doctor day: Weekday=%d Start=%v End=%v", day.Weekday, day.StartTime, day.EndTime)
				if day.Weekday == int(weekday) && !derefBool(day.IsDayOff) {
					slotMins := derefInt(day.SlotDurationMinutes)
					if slotMins == 0 {
						slotMins = 30
					}
					start := time.Date(date.Year(), date.Month(), date.Day(), day.StartTime.Hour(), day.StartTime.Minute(), 0, 0, time.Local)
					end := time.Date(date.Year(), date.Month(), date.Day(), day.EndTime.Hour(), day.EndTime.Minute(), 0, 0, time.Local)
					log.Printf("Generating slots for weekday=%d from %v to %v", weekday, start, end)
					slots = append(slots, generateSlots(start, end, slotMins, busy[dateBusyStr])...)
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

	sort.Slice(tempSlots, func(i, j int) bool {
		return tempSlots[i].Date.Before(tempSlots[j].Date)
	})

	var result []model.ScheduleEntry
	for _, slot := range tempSlots {
		result = append(result, model.ScheduleEntry{
			Label: slot.DateLabel,
			Slots: slot.Slots,
		})
	}

	//log.Println(result)
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

func (s *PatientService) GetUpcomingAppointments(ctx context.Context, token string) ([]model.UpcomingAppointment, error) {
	user, err := s.AuthClient.Client.GetPatient(ctx, &authpb.GetPatientRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	apps, err := s.StorageClient.Client.GetAppointmentsByUserID(ctx, &storagepb.GetByIDRequest{Id: user.Patient.UserId})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить предстоящие записи: %w", err)
	}
	upcoming := make([]model.UpcomingAppointment, 0)
	allSpecs, err := s.StorageClient.Client.GetAllSpecs(ctx, &storagepb.EmptyRequest{})
	specsMap := make(map[int]string)
	for _, spec := range allSpecs.Specs {
		specsMap[int(spec.Id)] = spec.Name
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось получить все специальности")
	}
	now := time.Now()
	for _, app := range apps.Appointment {
		if app.Status == "cancelled" {
			continue
		}
		if app.Date.AsTime().Before(now) {
			continue
		}
		if app.Time.AsTime().AddDate(now.Year(), int(now.Month()), now.Day()).Before(now) {
			continue
		}
		doc, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: app.DoctorId})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить врача для отображения записи: %w", err)
		}
		specs, err := s.StorageClient.Client.GetSpecsByDoctorID(ctx, &storagepb.GetByIDRequest{Id: doc.Doctor.UserId})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить специальности врача для отображения записи: %w", err)
		}
		docName := fmt.Sprintf("%s %s.%s.", doc.Doctor.SecondName, getAndCapitalizeFirstLetter(doc.Doctor.FirstName), getAndCapitalizeFirstLetter(doc.Doctor.Surname))
		upcomingDate := app.Date.AsTime().Format("02.01.2006")
		upcomingTime := app.Time.AsTime().Format("15:04")
		var specsString string
		for _, spec := range specs.SpecId {
			specName := fmt.Sprintf("%s, ", specsMap[int(spec)])
			specsString += specName
		}
		specsString = specsString[:len(specsString)-2]
		upcomingAppointment := model.UpcomingAppointment{
			ID:        model.AppointmentID(app.Id),
			Date:      upcomingDate,
			Time:      upcomingTime,
			DoctorID:  model.UserID(app.DoctorId),
			Doctor:    docName,
			Specialty: specsString,
		}
		upcoming = append(upcoming, upcomingAppointment)
	}
	return upcoming, nil
}

func (s *PatientService) UpdateAppointment(ctx context.Context, appointment model.Appointment) error {
	currApp, err := s.StorageClient.Client.GetAppointmentByID(ctx, &storagepb.GetByIDRequest{Id: int32(appointment.ID)})
	if err != nil {
		return fmt.Errorf("не удалось получить запись: %w", err)
	}
	updateApp := &storagepb.Appointment{
		Id:        int32(appointment.ID),
		Date:      timestamppb.New(appointment.Date),
		Time:      timestamppb.New(appointment.Time),
		Status:    currApp.Appointment.Status,
		UpdatedAt: timestamppb.New(time.Now()),
	}
	_, err = s.StorageClient.Client.UpdateAppointment(ctx, &storagepb.UpdateAppointmentRequest{Appointment: updateApp})
	if err != nil {
		return fmt.Errorf("не удалось обновить запись: %w", err)
	}
	return nil
}

func (s *PatientService) CancelAppointment(ctx context.Context, id model.AppointmentID) error {
	currApp, err := s.StorageClient.Client.GetAppointmentByID(ctx, &storagepb.GetByIDRequest{Id: int32(id)})
	if err != nil {
		return fmt.Errorf("не удалось получить запись: %w", err)
	}
	updateApp := &storagepb.Appointment{
		Id:        int32(id),
		Date:      currApp.Appointment.Date,
		Time:      currApp.Appointment.Time,
		Status:    "cancelled",
		UpdatedAt: timestamppb.New(time.Now()),
	}
	_, err = s.StorageClient.Client.UpdateAppointment(ctx, &storagepb.UpdateAppointmentRequest{Appointment: updateApp})
	if err != nil {
		return fmt.Errorf("не удалось отменить запись: %w", err)
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

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func getAndCapitalizeFirstLetter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	return string(unicode.ToUpper(runes[0]))
}

func generateSlots(start, end time.Time, slotMins int, busy map[string]bool) []string {
	if !start.Before(end) {
		log.Printf("WARNING: start >= end: %v >= %v", start, end)
		return nil
	}
	var result []string
	for t := start; t.Before(end); t = t.Add(time.Duration(slotMins) * time.Minute) {
		if t.Before(time.Now()) {
			continue
		}
		timeStr := t.Format("15:04")
		if busy != nil && busy[timeStr] {
			continue
		}
		result = append(result, timeStr)
	}
	return result
}
