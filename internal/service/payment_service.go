package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	patientpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/patient"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/repository"
	"github.com/razorpay/razorpay-go"
)

type PaymentService interface {
	// CreateAppointmentFeePayment(ctx context.Context, req *payment.CreateAppointmentFeePaymentRequest) (*payment.CreateAppointmentFeePaymentResponse, error)
	CreateRozorOrderId(reqData domain.Payment) (string, string, error)
	ProcessPayment(payment domain.Payment) error
}

type paymentService struct {
	repo           repository.PaymentRepository
	razorpayClient *razorpay.Client
	PatientClient  patientpb.PatientServiceClient
}

func NewPaymentService(repo repository.PaymentRepository, razorpayClient *razorpay.Client, patientClient patientpb.PatientServiceClient) PaymentService {
	return &paymentService{
		repo:           repo,
		razorpayClient: razorpayClient,
		PatientClient:  patientClient,
	}
}
func (p *paymentService) CreateRozorOrderId(reqData domain.Payment) (string, string, error) {
	orderParams := map[string]interface{}{
		"amount":   reqData.Amount * 100,
		"currency": "INR",
		"receipt":  reqData.PatientID,
		"notes": map[string]interface{}{
			"patientId": reqData.PatientID, // Add patientId to notes
		},
	}
	order, err := p.razorpayClient.Order.Create(orderParams, nil)
	if err != nil {
		log.Println("Order id creation failed ==", err)
		return "", "", errors.New("payment not initiated : " + err.Error())
	}
	razorId, _ := order["id"].(string)
	paymentUrl := "http://localhost:8080/api/v1/payment?orderId=" + razorId

	reqData.OrderID = razorId
	if err := p.repo.CreatePayment(&reqData); err != nil {
		return "Failed to store payment data", "", err
	}

	return razorId, paymentUrl, nil

}

func (s *paymentService) ProcessPayment(payment domain.Payment) error {
	// if err := RazorPaymentVerification(signature, orderID, paymentID); err != nil {
	// 	log.Println("Payment failed : Payment Signature not valid")
	// 	return err
	// }
	err := s.repo.UpdatePaymentStatus(payment)
	if err != nil {
		log.Printf("Error updating payment status: %v", err)
		return err
	}
	fmt.Println("patient_id", payment.PatientID)
	resp, err := s.PatientClient.GetProfile(context.Background(), &patientpb.GetProfileRequest{
		PatientId: payment.PatientID,
	})
	if err != nil {
		log.Println("Failed to get patient profile")
		return err
	}
	fmt.Println("email==", resp.Email)
	if payment.Status == "captured" {
		err = config.HandlePaymentCompletion(payment.OrderID, payment.PatientID, resp.Email, payment.Amount)
		if err != nil {
			return err
		}
	}
	log.Printf("Payment processed successfully for OrderID: %s", payment.OrderID)
	return nil
}

// func RazorPaymentVerification(sign, orderId, paymentId string) error {
// 	signature := sign
// 	secret := os.Getenv("RAZORPAY_KEY_SECRET")
// 	data := orderId + "|" + paymentId
// 	h := hmac.New(sha256.New, []byte(secret))
// 	_, err := h.Write([]byte(data))
// 	if err != nil {
// 		panic(err)
// 	}
// 	sha := hex.EncodeToString(h.Sum(nil))
// 	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
// 		return errors.New("payment signature not valid")
// 	} else {
// 		return nil
// 	}
// }

// func (s *paymentService) PaymentCallBack(ctx context.Context, req *payment.CreateAppointmentFeePaymentRequest) (*payment.CreateAppointmentFeePaymentResponse, error) {
// 	// Create Razorpay order
// 	orderData := map[string]interface{}{
// 		"amount":   int(req.Amount * 100), // amount in paise
// 		"currency": "INR",
// 		"receipt":  fmt.Sprintf("receipt_%s", req.PatientId),
// 	}

// 	order, err := s.razorpayClient.Order.Create(orderData, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create Razorpay order: %v", err)
// 	}

// 	// Create payment record in database
// 	paymentRecord := &domain.Payment{
// 		PatientID:        req.PatientId,
// 		SpecializationID: req.SpecializationId,
// 		Amount:           req.Amount,
// 		Status:           "pending",
// 		PaymentID:        "", // Will be updated upon successful payment
// 		OrderID:          order["id"].(string),
// 		Type:             req.Type,
// 	}

// 	if err := s.repo.CreatePayment(paymentRecord); err != nil {
// 		return nil, fmt.Errorf("failed to create payment record: %v", err)
// 	}

// 	// Return response with order ID
// 	return &payment.CreateAppointmentFeePaymentResponse{
// 		PaymentId:  paymentRecord.PaymentID,
// 		Status:     paymentRecord.Status,
// 		Message:    "Order created successfully",
// 		PaymentUrl: paymentRecord.OrderID,
// 	}, nil
// }
