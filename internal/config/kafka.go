package config

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
}


type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	PatientID string  `json:"patient_id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
}
// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) ProducePaymentEvent(topic string, event PaymentEvent) error {
	// Serialize the event to JSON
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	deliveryChan := make(chan kafka.Event)
	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message, // Send serialized JSON message
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
	}

	log.Printf("Message delivered to topic %s [%d] at offset %v\n", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	close(deliveryChan)

	return nil
}

// HandlePaymentCompletion handles payment and produces an event
func HandlePaymentCompletion(paymentID, patientID, email string, amount float64) error {
	kafkaProducer, err := NewKafkaProducer("localhost:9092") // Use your Kafka broker address
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Create a payment event struct
	paymentEvent := PaymentEvent{
		PaymentID: paymentID,
		PatientID: patientID,
		Email:     email,
		Amount:    amount,
		Date:      time.Now().Format(time.RFC3339), // Use current timestamp
	}

	// Produce the event to Kafka
	err = kafkaProducer.ProducePaymentEvent("payment_topic", paymentEvent)
	if err != nil {
		return fmt.Errorf("failed to produce payment event: %w", err)
	}

	log.Println("Payment event produced successfully")
	return nil
}
