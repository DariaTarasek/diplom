package model

import "time"

type DoctorSchedule struct {
	ID                  int        `db:"id"`
	DoctorID            UserID     `db:"doctor_id"`
	Weekday             int        `db:"weekday"`
	StartTime           *time.Time `db:"start_time"`
	EndTime             *time.Time `db:"end_time"`
	SlotDurationMinutes *int       `db:"slot_duration_minutes"`
	IsDayOff            *bool      `db:"is_day_off"`
}

type ClinicSchedule struct {
	ID                  int        `db:"id"`
	Weekday             int        `db:"weekday"`
	StartTime           *time.Time `db:"start_time"`
	EndTime             *time.Time `db:"end_time"`
	SlotDurationMinutes *int       `db:"slot_duration_minutes"`
	IsDayOff            *bool      `db:"is_day_off"`
}
