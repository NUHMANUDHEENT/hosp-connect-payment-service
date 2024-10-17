package di

import (
	"log"
	"net"

	patientpb "github.com/NUHMANUDHEENT/hosp-connect-pb/proto/patient"
	"github.com/NUHMANUDHEENT/hosp-connect-pb/proto/payment"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/handler"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/repository"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/service"
	"github.com/razorpay/razorpay-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCSetup initializes the gRPC server and registers the services
func GRPCSetup(port string, razorpayClient *razorpay.Client) (net.Listener, *grpc.Server) {
	// Create a TCP listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	db := config.InitDatabase()
	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(db)

	patientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Doctor client calling
	if err != nil {
		log.Fatalf("Failed to connect to patient service: %v", err)
	}
	patientClient := patientpb.NewPatientServiceClient(patientConn)
	// Initialize services
	paymentService := service.NewPaymentService(paymentRepo, razorpayClient, patientClient)

	// Initialize handlers
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the PaymentService with the gRPC server
	payment.RegisterPaymentServiceServer(server, paymentHandler)

	// Enable server reflection (optional, useful for testing with tools like grpcurl)
	reflection.Register(server)
	log.Printf("Payment gRPC service is running on port %s", port)
	return listener, server
}
