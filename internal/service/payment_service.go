package service

import (
	"context"
	"errors"

	appointmentpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/appointment"
	patientpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/patient"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/utils"
	"github.com/razorpay/razorpay-go"
	"github.com/sirupsen/logrus"
)

type PaymentService interface {
	// CreateAppointmentFeePayment(ctx context.Context, req *payment.CreateAppointmentFeePaymentRequest) (*payment.CreateAppointmentFeePaymentResponse, error)
	CreateRozorOrderId(reqData domain.Payment) (string, string, error)
	ProcessPayment(payment domain.Payment) error
	GetTotalRevenue(param string) (float64, error)
}

type paymentService struct {
	repo              repository.PaymentRepository
	razorpayClient    *razorpay.Client
	PatientClient     patientpb.PatientServiceClient
	AppointmentClient appointmentpb.AppointmentServiceClient
	Logger            *logrus.Logger
}

func NewPaymentService(repo repository.PaymentRepository, razorpayClient *razorpay.Client, patientClient patientpb.PatientServiceClient, appointmentClient appointmentpb.AppointmentServiceClient, logger *logrus.Logger) PaymentService {
	return &paymentService{
		repo:              repo,
		razorpayClient:    razorpayClient,
		PatientClient:     patientClient,
		AppointmentClient: appointmentClient,
		Logger:            logger,
	}
}
func (p *paymentService) CreateRozorOrderId(reqData domain.Payment) (string, string, error) {
	p.Logger.WithFields(logrus.Fields{
		"Function":  "CreateRozorOrderId",
		"PatientID": reqData.PatientID,
		"Amount":    reqData.Amount,
	}).Info("Initiating Razorpay order creation")

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
		p.Logger.WithError(err).Error("Failed to create Razorpay order")
		return "", "", errors.New("payment not initiated : " + err.Error())
	}
	razorId, _ := order["id"].(string)
	paymentUrl := "http://localhost:8080/api/v1/payment?orderId=" + razorId

	reqData.OrderID = razorId
	if err := p.repo.CreatePayment(&reqData); err != nil {
		p.Logger.WithError(err).Error("Failed to store payment data")
		return "Failed to store payment data", "", err
	}

	p.Logger.WithFields(logrus.Fields{
		"OrderID": razorId,
		"URL":     paymentUrl,
	}).Info("Razorpay order created successfully")
	return razorId, paymentUrl, nil
}

func (s *paymentService) ProcessPayment(payment domain.Payment) error {
	s.Logger.WithFields(logrus.Fields{
		"Function": "ProcessPayment",
		"OrderID":  payment.OrderID,
		"Status":   payment.Status,
	}).Info("Processing payment")

	err := s.repo.UpdatePaymentStatus(payment)
	if err != nil {
		s.Logger.WithError(err).Error("Error updating payment status")
		return err
	}

	resp, err := s.PatientClient.GetProfile(context.Background(), &patientpb.GetProfileRequest{
		PatientId: payment.PatientID,
	})
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch patient profile")
		return err
	}

	s.Logger.WithFields(logrus.Fields{
		"Email": resp.Email,
	}).Info("Fetched patient profile for payment notification")

	if payment.Status == "captured" {
		s.Logger.Info("Payment captured, fetching appointment details")

		appointment, err := s.AppointmentClient.GetAppointmentDetails(context.Background(), &appointmentpb.GetAppointmentDetailsRequest{
			OrderId: payment.OrderID,
		})
		if err != nil {
			s.Logger.WithError(err).Error("Failed to fetch appointment details")
			return err
		}

		err = utils.HandleAppointmentNotification(payment.OrderID, payment.PatientID, resp.Email, payment.Amount, appointment.AppointmentTime.AsTime())
		if err != nil {
			s.Logger.WithError(err).Error("Failed to send payment notification")
			return err
		}
	}

	s.Logger.WithFields(logrus.Fields{
		"OrderID": payment.OrderID,
		"Status":  payment.Status,
	}).Info("Payment processed successfully")
	return nil
}

func (s *paymentService) GetTotalRevenue(param string) (float64, error) {
	s.Logger.WithFields(logrus.Fields{
		"Function": "GetTotalRevenue",
		"Param":    param,
	}).Info("Fetching total revenue")

	revenue, err := s.repo.GetTotalRevenue(param)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch total revenue")
		return 0, err
	}

	s.Logger.WithFields(logrus.Fields{
		"TotalRevenue": revenue,
	}).Info("Total revenue fetched successfully")
	return revenue, nil
}
