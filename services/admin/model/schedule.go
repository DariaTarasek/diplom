package model

import "time"

type ClinicWeeklySchedule struct {
	ID                  int
	Weekday             int
	StartTime           time.Time
	EndTime             time.Time
	SlotDurationMinutes int
	IsDayOff            bool
}

type DoctorWeeklySchedule struct {
	ID                  int
	DoctorID            int
	Weekday             int
	StartTime           time.Time
	EndTime             time.Time
	SlotDurationMinutes int
	IsDayOff            bool
}

type ClinicDailyOverride struct {
	ID                  int
	Date                time.Time
	StartTime           time.Time
	EndTime             time.Time
	SlotDurationMinutes int
	IsDayOff            bool
}

type DoctorDailyOverride struct {
	ID                  int
	DoctorId            int
	Date                time.Time
	StartTime           time.Time
	EndTime             time.Time
	SlotDurationMinutes int
	IsDayOff            bool
}
