package schema

import (
	"time"

	"gorm.io/gorm"
)

type SalaryInfo struct {
	gorm.Model
	UserID        uint      `json:"userId"` // Foreign key for User
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	EffectiveDate time.Time `json:"effectiveDate"` // Date from which the current salary is effective
}

type SalaryPayment struct {
	gorm.Model
	UserID        uint      `json:"userId"`        // Foreign key for User
	SalaryInfoID  uint      `json:"salaryInfoId"`  // Foreign key for SalaryInfo
	PaymentDate   string    `json:"paymentDate"`   // Date of payment
	Amount        float64   `json:"amount"`        // Total amount paid
	PaymentMethod string    `json:"paymentMethod"` // Method of payment
	Status        string    `json:"status"`        // Payment status (e.g., Completed, Pending)
	PeriodStart   time.Time `json:"periodStart"`   // Start of the pay period
	PeriodEnd     time.Time `json:"periodEnd"`     // End of the pay period
	Deductions    float64   `json:"deductions"`    // Deductions, if any
	Bonuses       float64   `json:"bonuses"`       // Bonuses, if any
	Notes         string    `json:"notes"`         // Additional notes
}

func (salaryInfo *SalaryInfo) AfterCreate(tx *gorm.DB) (err error) {
	// Update User's SalaryInfoID field
	err = tx.Model(&User{}).Where("id = ?", salaryInfo.UserID).Update("current_salary_info_id", salaryInfo.ID).Error
	if err != nil {
		return err
	}
	return nil
}
