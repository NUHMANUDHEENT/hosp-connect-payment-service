package main

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/di"
)

func main() {
	config.LoadEnv()

	port := os.Getenv("PAYMENT_PORT")
	if port == "" {
		log.Fatalf("PAYMENT_PORT not set")
	}
	listener, server := config.GRPCSetup(port, di.RazorClientSetUp())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
