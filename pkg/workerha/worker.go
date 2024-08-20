package workerha

import (
	"context"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"server-management/pkg/elasticha"
	"server-management/pkg/kafkaha"
	"server-management/pkg/postgresha"
)

// Worker struct holds dependencies for the worker
type Worker struct {
	KafkaConsumer kafkaha.ConsumerClientInterface
	Elasticsearch *elasticha.Repository
	PostgresRepo  *postgresha.Repository[Server]
}

// NewWorker initializes a new Worker
func NewWorker(kafkaConsumer kafkaha.ConsumerClientInterface, es *elasticha.Repository, pgRepo *postgresha.Repository[Server]) *Worker {
	return &Worker{
		KafkaConsumer: kafkaConsumer,
		Elasticsearch: es,
		PostgresRepo:  pgRepo,
	}
}

// Server represents the server model from your database
type Server struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by"`
	DeletedBy string         `json:"deleted_by"`
	Name      string         `json:"name"`
	IPv4      string         `json:"ipv4"`
}

// pingServer performs a ping to the server's IPv4 address
func (w *Worker) pingServer(ip string) string {
	cmd := exec.Command("ping", "-c", "1", ip)
	err := cmd.Run()
	if err != nil {
		return "offline"
	}
	return "online"
}

// processJob processes a single job from Kafka
// processJob processes a single job from Kafka
func (w *Worker) processJob(ctx context.Context, job kafkaha.ServerStatus, wg *sync.WaitGroup) {
	defer wg.Done() // This will now work correctly

	// Check server status
	status := w.pingServer(job.ServerID)

	// Create ServerDocument for Elasticsearch
	serverDoc := elasticha.ServerDocument{
		ID:     uuid.New(),
		Name:   job.ServerID,
		Ipv4:   job.ServerID, // Assuming ServerID is the IP, adapt as necessary
		Status: status,
		At:     time.Now(),
		In:     1, // Example data, adapt as necessary
	}

	// Save to Elasticsearch
	_, err := w.Elasticsearch.CreateOne(serverDoc)
	if err != nil {
		return err
	}

	// Update the status in PostgreSQL
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if _, err := w.PostgresRepo.UpdateOneById(ctx, serverDoc.ID, updateData); err != nil {
		log.Printf("Error updating server status in PostgreSQL: %v", err)
	}
}

// Start begins the worker loop
func (w *Worker) Start(ctx context.Context) {
	for {
		// Read health check jobs from Kafka
		jobs, err := w.KafkaConsumer.ReadStruct(ctx)
		if err != nil {
			log.Printf("Error reading from Kafka: %v", err)
			continue
		}

		// Use a WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup

		// Process each job concurrently
		for _, job := range jobs {
			wg.Add(1) // Increment the WaitGroup counter
			go w.processJob(ctx, job, &wg)
		}

		wg.Wait() // Block until all goroutines have finished

		// Sleep before checking Kafka for more jobs (adjust the interval as needed)
		time.Sleep(5 * time.Minute)
	}
}
