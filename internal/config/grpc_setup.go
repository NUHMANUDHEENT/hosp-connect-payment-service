package config

import (
	"log"
	"net"
	"os"

	appointmentpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/appointment"
	patientpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/patient"
	"github.com/NUHMANUDHEENT/hosp-connect-pb/proto/payment"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/handler"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/service"
	"github.com/nuhmanudheent/hosp-connect-payment-service/logs"
	"github.com/razorpay/razorpay-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func GRPCSetup(port string, razorpayClient *razorpay.Client) (net.Listener, *grpc.Server) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	db := InitDatabase()

	logger := logs.NewLogger()
	paymentRepo := repository.NewPaymentRepository(db)

	patientConn, err := grpc.NewClient(os.Getenv("USER_GRPC_SERVER"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to patient service: %v", err)
	}
	appointmentConn, err := grpc.NewClient(os.Getenv("APPT_GRPC_SERVER"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to appointment service: %v", err)
	}
	patientClient := patientpb.NewPatientServiceClient(patientConn)
	appointmentClient := appointmentpb.NewAppointmentServiceClient(appointmentConn)

	paymentService := service.NewPaymentService(paymentRepo, razorpayClient, patientClient, appointmentClient, logger)

	paymentHandler := handler.NewPaymentHandler(paymentService)
	
	server := grpc.NewServer()

	payment.RegisterPaymentServiceServer(server, paymentHandler)


	reflection.Register(server)
	log.Printf("Payment gRPC service is running on port %s", port)
	return listener, server
}
