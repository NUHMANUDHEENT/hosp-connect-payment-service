package main

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/di" // Import Razorpay SDK
)

func main() {
	di.LoadEnv()

	port := os.Getenv("PAYMENT_PORT")
	if port == "" {
		log.Fatalf("PAYMENT_PORT not set")
	}
	// Call GRPCSetup and pass the razorpayClient
	listener, server := di.GRPCSetup(port, di.RazorClientSetUp())

	// Start the gRPC server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
	
}
