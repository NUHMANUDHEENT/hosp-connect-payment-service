package domain

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PatientID     string `gorm:"not null"`
	AppointmentId int
	Amount        float64 `gorm:"not null"`
	Status        string  `gorm:"not null"`
	PaymentID     string
	OrderID       string `gorm:"uniqueIndex"`
	Type          string
}
