package schema

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	UserID   uint      `json:"user_id"   gorm:"not null"`
	Date     time.Time `json:"date"      gorm:"type:date;not null"`
	ClockIn  time.Time `json:"clock_in"  gorm:"type:time"`
	ClockOut time.Time `json:"clock_out" gorm:"type:time"`
	Notes    string    `json:"notes"     gorm:"type:text"`
}
