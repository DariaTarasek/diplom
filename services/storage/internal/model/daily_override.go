package model

import "time"

type DoctorDailyOverride struct {
	ID                  int        `db:"id"`
	DoctorID            UserID     `db:"doctor_id"`
	Date                time.Time  `db:"date"`
	StartTime           *time.Time `db:"start_time"`
	EndTime             *time.Time `db:"end_time"`
	SlotDurationMinutes *int       `db:"slot_duration_minutes"`
	IsDayOff            *bool      `db:"is_day_off"`
}

type ClinicDailyOverride struct {
	ID                  int        `db:"id"`
	Date                time.Time  `db:"date"`
	StartTime           *time.Time `db:"start_time"`
	EndTime             *time.Time `db:"end_time"`
	SlotDurationMinutes *int       `db:"slot_duration_minutes"`
	IsDayOff            *bool      `db:"is_day_off"`
}
