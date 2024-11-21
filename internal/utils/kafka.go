package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	PatientID string  `json:"patient_id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
}

func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        "payment_topic",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireOne),
	})

	return &KafkaProducer{writer: writer}, nil
}

func (kp *KafkaProducer) ProducePaymentEvent(event PaymentEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.PaymentID),
		Value: message,
	}

	// Send the message
	err = kp.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	log.Printf("Message successfully sent to Kafka topic: %s", kp.writer.Topic)
	return nil
}


func HandleAppointmentNotification(paymentID, patientID, email string, amount float64, datetime time.Time) error {
	kafkaProducer, err := NewKafkaProducer(os.Getenv("KAFKA_BROKER"))
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	paymentEvent := PaymentEvent{
		PaymentID: paymentID,
		PatientID: patientID,
		Email:     email,
		Amount:    amount,
		Date:      time.Now().Format(time.RFC3339),
	}


	err = kafkaProducer.ProducePaymentEvent(paymentEvent)
	if err != nil {
		return fmt.Errorf("failed to produce payment event: %w", err)
	}

	log.Println("Payment event produced successfully")
	return nil
}
func EnsureTopicExists(broker, topic string) error {
    conn, err := kafka.Dial("tcp", broker)
    if err != nil {
        return fmt.Errorf("failed to connect to Kafka broker: %w", err)
    }
    defer conn.Close()

    topics, err := conn.ReadPartitions()
    if err != nil {
        return fmt.Errorf("failed to read partitions: %w", err)
    }

    for _, t := range topics {
        if t.Topic == topic {
            return nil
        }
    }

    // Create the topic if it doesn't exist
    err = conn.CreateTopics(kafka.TopicConfig{
        Topic:             topic,
        NumPartitions:     -1,
        ReplicationFactor: -1,
    })
    if err != nil {
        return fmt.Errorf("failed to create topic: %w", err)
    }

    return nil
}
