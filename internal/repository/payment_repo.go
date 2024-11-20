package repository

import (
	"log"

	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *domain.Payment) error
	GetPaymentByID(paymentID string) (*domain.Payment, error)
	UpdatePaymentStatus(payment domain.Payment) error
	GetTotalRevenue(param string) (float64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) CreatePayment(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetPaymentByID(paymentID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("payment_id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
func (r *paymentRepository) UpdatePaymentStatus(payment domain.Payment) error {

	if err := r.db.Where("order_id =?", payment.OrderID).Model(&payment).Updates(payment).Error; err != nil {
		return err
	}
	log.Println("payment details updated")
	return nil
}
func (r *paymentRepository) GetTotalRevenue(param string) (float64, error) {
	var totalRevenue float64

	query := r.db.Model(&domain.Payment{}).Select("SUM(amount)")

	switch param {
	case "day":
		query = query.Where("DATE_TRUNC('day', created_at) = DATE_TRUNC('day', CURRENT_TIMESTAMP)")
	case "week":
		query = query.Where("DATE_TRUNC('week', created_at) = DATE_TRUNC('week', CURRENT_TIMESTAMP)")
	case "month":
		query = query.Where("DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_TIMESTAMP)")
	}

	if err := query.Scan(&totalRevenue).Error; err != nil {
		return 0, err
	}

	return totalRevenue, nil
}
