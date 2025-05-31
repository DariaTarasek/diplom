package model

type ClinicWeeklySchedule struct {
	ID                  int    `json:"id"`
	Weekday             int    `json:"day"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	SlotDurationMinutes int    `json:"slot_minutes"`
	IsDayOff            bool   `json:"is_day_off"`
}

type DoctorWeeklySchedule struct {
	ID                  int    `json:"id"`
	DoctorID            int    `json:"selectedDoctor"`
	Weekday             int    `json:"day"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	SlotDurationMinutes int    `json:"slot_minutes"`
	IsDayOff            bool   `json:"is_day_off"`
}

type ClinicDailyOverride struct {
	ID        int    `json:"id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsDayOff  string `json:"type"`
}

type DoctorDailyOverride struct {
	ID        int    `json:"id"`
	DoctorId  int    `json:"doctor_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsDayOff  string `json:"type"`
}
