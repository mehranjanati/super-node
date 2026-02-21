package podman

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"go.uber.org/fx"
)

// ManageInfrastructure manages the lifecycle of the infrastructure containers.
func ManageInfrastructure(lc fx.Lifecycle) {
	// Determine socket path based on OS and availability
	socketPath := "/run/podman/podman.sock"

	// Check for common socket locations
	commonSockets := []string{
		"/var/run/docker.sock",
		"/run/docker.sock",
		"/run/podman/podman.sock",
	}

	for _, s := range commonSockets {
		if _, err := os.Stat(s); err == nil {
			socketPath = s
			break
		}
	}

	if runtime.GOOS == "darwin" {
		// On macOS, explicitly prefer default Docker socket if not found above
		if _, err := os.Stat("/var/run/docker.sock"); err == nil {
			socketPath = "/var/run/docker.sock"
		}
	}

	// Create a custom HTTP client that talks over the Unix socket
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
		Timeout: 5 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Checking infrastructure health via socket: %s", socketPath)

			// List of critical containers to check
			// Note: Names match the docker-compose/podman generated names
			criticalContainers := []string{"podman-tidb-1", "podman-hasura-1", "podman-redpanda-1"}

			for _, name := range criticalContainers {
				// We use the Docker-compatible API which Podman supports
				// GET /containers/{name}/json
				resp, err := client.Get(fmt.Sprintf("http://localhost/containers/%s/json", name))
				if err != nil {
					log.Printf("Warning: Failed to check container %s (socket might be unreachable): %v", name, err)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusNotFound {
					log.Printf("Warning: Container %s NOT FOUND", name)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					log.Printf("Warning: Unexpected status for container %s: %d", name, resp.StatusCode)
					continue
				}

				var data struct {
					State struct {
						Running bool   `json:"Running"`
						Status  string `json:"Status"`
					} `json:"State"`
				}

				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					log.Printf("Warning: Failed to decode response for %s: %v", name, err)
					continue
				}

				if data.State.Running {
					log.Printf("Container %s is RUNNING", name)
				} else {
					log.Printf("Warning: Container %s is NOT running (State: %s)", name, data.State.Status)
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
