package model

type ClinicWeeklySchedule struct {
	ID                  int    `json:"id"`
	Weekday             int    `json:"day"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	SlotDurationMinutes int    `json:"slot_minutes"`
	IsDayOff            bool   `json:"is_day_off"`
}
