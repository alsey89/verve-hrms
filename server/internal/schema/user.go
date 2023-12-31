package schema

import (
	"time"

	"gorm.io/gorm"
)

// User model ---------------------------------------------------------------
type User struct {
	//* account information
	gorm.Model            // This includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	Email      string     `json:"email"      gorm:"uniqueIndex"`
	Password   string     `json:"-"          gorm:"type:varchar(100)"` //* Password is not returned in JSON
	IsAdmin    bool       `json:"isAdmin"    gorm:"default:false"`
	AvatarURL  *string    `json:"avatarUrl"  gorm:"type:text"`
	IsActive   bool       `json:"isActive"   gorm:"default:true"`
	IsVerified bool       `json:"isVerified" gorm:"default:false"`
	LastLogin  *time.Time `json:"lastLogin"  gorm:"default:null"`

	//* personal information
	FirstName        string            `json:"firstName"`
	MiddleName       *string           `json:"middleName"`
	LastName         string            `json:"lastName"`
	Nickname         *string           `json:"nickname"`
	ContactInfo      *ContactInfo      `json:"contactInfo"`
	EmergencyContact *EmergencyContact `json:"emergencyContact"`

	//*job information
	Title      string   `json:"title"      gorm:"default:null"` //! DENORMALIZED DATA WITH SYNC CALLBACK (see job.go)
	Department string   `json:"department" gorm:"default:null"` //! DENORMALIZED DATA WITH SYNC CALLBACK (see job.go)
	JobInfo    *JobInfo `json:"jobInfo"`

	//* salary information
	SalaryInfoID  uint            `json:"salaryInfoId"`
	SalaryPayment []SalaryPayment `json:"salaryPayment"`
}

type ContactInfo struct {
	gorm.Model
	UserID     uint   `json:"userId"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
}

type EmergencyContact struct {
	gorm.Model
	UserID     uint    `json:"userId"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Relation   string  `json:"relation"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
}

type Job struct {
	gorm.Model
	UserID     uint   `json:"userId"`
	JobTitle   string `json:"jobTitle"`
	Department string `json:"department"`
	Location   string `json:"location"`
	Manager    string `json:"manager"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}
