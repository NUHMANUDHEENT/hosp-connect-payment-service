package di

import (
	"log"
	"os"

	"github.com/razorpay/razorpay-go"
)

func RazorClientSetUp() *razorpay.Client {
	razorpayKeyID := os.Getenv("RAZORPAY_KEY_ID")
	razorpayKeySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	if razorpayKeyID == "" || razorpayKeySecret == "" {
		log.Fatalf("Razorpay credentials not set")
	}

	razorpayClient := razorpay.NewClient(razorpayKeyID, razorpayKeySecret)
	return razorpayClient
}
