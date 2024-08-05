package kafkaha

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

type ServerStatus struct {
	ServerID string `json:"server_id"`
	Status   string `json:"status"`
}

type ConsumerClientInterface interface {
	ReadMessage(echo.Context) (kafka.Message, error)
	CloseReader() error
	ReadStruct(echo.Context) ([]string, error)
}

type ConsumerKafka struct {
	k *kafka.Reader
}

func (c *ConsumerKafka) ReadMessage(ctx echo.Context) (kafka.Message, error) {
	return c.k.ReadMessage(ctx.Request().Context())
}

func (c *ConsumerKafka) CloseReader() error {
	return c.k.Close()
}

func (c *ConsumerKafka) ReadStruct(ctx echo.Context) ([]string, error) {
	msg, err := c.k.ReadMessage(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	var result []string
	if err := json.Unmarshal(msg.Value, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type ConfigConsumer struct {
	Broker      []string
	Topic       string
	GroupID     string
	Logger      kafka.Logger
	ErrorLogger kafka.Logger
}

func NewConsumer(c ConfigConsumer) *ConsumerKafka {
	return &ConsumerKafka{
		k: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: c.Broker,
				GroupID: c.GroupID,
				Topic:   c.Topic,
			},
		),
	}
}

type ProducerClientInterface interface {
	WriteMessage(echo.Context, ...kafka.Message) error
	WriteStruct(echo.Context, []ServerStatus) error
}

type ProducerKafka struct {
	k *kafka.Writer
}

func (c *ProducerKafka) WriteMessage(ctx echo.Context, msgs ...kafka.Message) error {
	return c.k.WriteMessages(ctx.Request().Context(), msgs...)
}

func (c *ProducerKafka) WriteStruct(ctx echo.Context, data []ServerStatus) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Value: jsonData,
	}
	return c.k.WriteMessages(ctx.Request().Context(), msg)
}

type ConfigProducer struct {
	Broker      []string
	Topic       string
	Balancer    kafka.Balancer
	Logger      kafka.Logger
	ErrorLogger kafka.Logger
}

func NewProducer(c ConfigProducer) *ProducerKafka {
	return &ProducerKafka{
		k: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers:  c.Broker,
				Topic:    c.Topic,
				Balancer: c.Balancer,
			},
		),
	}
}
