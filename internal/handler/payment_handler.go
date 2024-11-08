package handler

import (
	"context"
	"errors"
	"log"

	paymentpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/payment"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/domain"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/service"
	// "github.com/your_project/internal/service"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type PaymentServiceClient struct {
	paymentpb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentServiceClient {
	return &PaymentServiceClient{
		paymentService: service,
	}
}
func (h *PaymentServiceClient) CreateRazorOrderId(ctx context.Context, req *paymentpb.CreateRazorOrderIdRequest) (*paymentpb.CreateRazorOrderIdResponse, error) {
	log.Println("request with", req.PatientId)
	res, paymenturl, err := h.paymentService.CreateRozorOrderId(domain.Payment{
		PatientID: req.PatientId,
		Amount:    req.Amount,
		Type:      req.Type,
	})
	if err != nil {
		return &paymentpb.CreateRazorOrderIdResponse{
			Message:    err.Error(),
			Status:     "fail",
			StatusCode: 400,
		}, nil

	}
	return &paymentpb.CreateRazorOrderIdResponse{
		Message:    "Razorid created successfully",
		Status:     "success",
		StatusCode: 200,
		PaymentUrl: paymenturl,
		OrderId:    res,
	}, nil
}
func (h *PaymentServiceClient) PaymentCallback(ctx context.Context, req *paymentpb.PaymentCallBackRequest) (*paymentpb.PaymentCallBackResponse, error) {
	log.Printf("Received payment callback: OrderID: %s, PaymentID: %s, Status: %s", req.OrderId, req.PaymentId, req.Status)
	payment := domain.Payment{
		PatientID: req.PatientId,
		OrderID:   req.OrderId,
		PaymentID: req.PaymentId,
		Status:    req.Status,
		Amount:    req.Amount,
	}
	// Call the payment service to process and update the payment details
	err := h.paymentService.ProcessPayment(payment)
	if err != nil {
		log.Printf("Error processing payment: %v", err)
		return &paymentpb.PaymentCallBackResponse{
			Message: "Failed to process payment",
			Success: false,
		}, nil
	}

	// Successful processing
	return &paymentpb.PaymentCallBackResponse{
		Message: "Payment processed successfully",
		Success: true,
	}, nil
}
func (h *PaymentServiceClient) GetTotalRevenue(ctx context.Context, req *paymentpb.GetTotalRevenueRequest) (*paymentpb.GetTotalRevenueResponse, error) {
	revenue, err := h.paymentService.GetTotalRevenue(req.Param)
	if err != nil {
		return &paymentpb.GetTotalRevenueResponse{}, errors.New("failed to get total revenue")
	}
	return &paymentpb.GetTotalRevenueResponse{
		TotalRevenue: revenue,
	}, nil
}
