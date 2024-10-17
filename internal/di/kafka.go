package di

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil
}

// ProducePaymentEvent publishes the payment event to Kafka
func (kp *KafkaProducer) ProducePaymentEvent(topic, message string) error {
	deliveryChan := make(chan kafka.Event)
	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
	}

	fmt.Printf("Message delivered to topic %s [%d] at offset %v\n", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	close(deliveryChan)

	return nil
}
func HandlePaymentCompletion(paymentID string, amount float64) error {
    kafkaProducer, err := NewKafkaProducer("localhost:9092") // Use your Kafka broker address
    if err != nil {
        return fmt.Errorf("failed to create Kafka producer: %w", err)
    }

    message := fmt.Sprintf("Payment completed with ID: %s, Amount: %.2f", paymentID, amount)
    err = kafkaProducer.ProducePaymentEvent("payment_topic", message)
    if err != nil {
        return fmt.Errorf("failed to produce payment event: %w", err)
    }

    fmt.Println("Payment event produced successfully")
    return nil
}