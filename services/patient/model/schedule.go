package model

import "time"

type DoctorSchedule struct {
	ID                  int
	DoctorID            UserID
	Weekday             int
	StartTime           *time.Time
	EndTime             *time.Time
	SlotDurationMinutes *int
	IsDayOff            *bool
}
