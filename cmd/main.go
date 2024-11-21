package main

import (
	"log"
	"os"

	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/config"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/di"
	"github.com/nuhmanudheent/hosp-connect-payment-service/internal/utils"
)

func main() {
	config.LoadEnv()

	port := os.Getenv("PAYMENT_PORT")
	if port == "" {
		log.Fatalf("PAYMENT_PORT not set")
	}
	utils.EnsureTopicExists(os.Getenv("KAFKA_BROKER"),"payment_topic")
	listener, server := config.GRPCSetup(port, di.RazorClientSetUp())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
