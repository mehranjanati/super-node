package temporal

import (
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
)

func NewClient() (client.Client, error) {
	hostPort := os.Getenv("TEMPORAL_HOST_PORT")
	if hostPort == "" {
		hostPort = "temporal:7233"
	}

	var c client.Client
	var err error

	// Retry up to 30 times (30 seconds)
	for i := 0; i < 30; i++ {
		c, err = client.Dial(client.Options{
			HostPort: hostPort,
		})
		if err == nil {
			log.Printf("Successfully connected to Temporal at %s", hostPort)
			return c, nil
		}
		log.Printf("Failed to connect to Temporal at %s: %v. Retrying in 1s...", hostPort, err)
		time.Sleep(1 * time.Second)
	}

	log.Printf("Unable to create temporal client after retries: %v", err)
	return nil, err
}
