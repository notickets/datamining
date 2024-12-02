package datamining

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type Parser struct {
	KafkaWriter *kafka.Writer
	ProxyURL    string
	Client      *http.Client
}

func NewParser() (*Parser, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	proxyURL := os.Getenv("PROXY_URL")

	if kafkaBroker == "" || kafkaTopic == "" {
		return nil, errors.New("env variables PROXY_URL, KAFKA_BROKER, and KAFKA_TOPIC are required")
	}

	return &Parser{
		KafkaWriter: &kafka.Writer{
			Addr:         kafka.TCP(kafkaBroker),
			Topic:        kafkaTopic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100, // Increase batch size
			BatchTimeout: 10 * time.Millisecond,
		},
		ProxyURL: proxyURL,
		Client:   &http.Client{},
	}, nil
}

func (p *Parser) GetRequest(URL string) (string, error) {
	client := p.Client
	if p.ProxyURL != "" {
		proxy, err := url.Parse(p.ProxyURL)
		if err != nil {
			return "", err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	resp, err := client.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch data from " + URL + ", status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (p *Parser) DoRequest(method, URL string, body io.Reader) (string, error) {
	client := p.Client
	if p.ProxyURL != "" {
		proxy, err := url.Parse(p.ProxyURL)
		if err != nil {
			return "", err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch data from " + URL + ", status: " + resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

func (p *Parser) SendToKafka(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.Name),
		Value: data,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.KafkaWriter.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Message sent to Kafka: %s", event.Name)
	return nil
}

func (p *Parser) Close() error {
	if p.KafkaWriter == nil {
		return nil
	}
	return p.KafkaWriter.Close()
}
