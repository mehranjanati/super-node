package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Redpanda RedpandaConfig `mapstructure:"redpanda"`
	Hasura   HasuraConfig   `mapstructure:"hasura"`
	Agents   AgentsConfig   `mapstructure:"agents"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Rivet    RivetConfig    `mapstructure:"rivet"`
	OpenClaw OpenClawConfig `mapstructure:"openclaw"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	Benthos  BenthosConfig  `mapstructure:"benthos"`
	MCP      MCPConfig      `mapstructure:"mcp"`
	Matrix   MatrixConfig   `mapstructure:"matrix"`
	LiveKit  LiveKitConfig  `mapstructure:"livekit"`
	TiDB     TiDBConfig     `mapstructure:"tidb"`
	Role     string         `mapstructure:"role"` // monolith, api, worker, consumer
}

type TiDBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type MatrixConfig struct {
	HomeserverURL            string `mapstructure:"homeserver_url"`
	ServerName               string `mapstructure:"server_name"`
	RegistrationSharedSecret string `mapstructure:"registration_shared_secret"`
}

type LiveKitConfig struct {
	APIURL    string `mapstructure:"api_url"`
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`
}

type MCPConfig struct {
	Servers []MCPServerConfig `mapstructure:"servers"`
}

type MCPServerConfig struct {
	ID          string   `mapstructure:"id"`
	Name        string   `mapstructure:"name"`
	Type        string   `mapstructure:"type"` // "stdio" or "sse"
	Command     string   `mapstructure:"command"`
	Args        []string `mapstructure:"args"`
	URL         string   `mapstructure:"url"` // for sse
	Environment []string `mapstructure:"env"`
}

type BenthosConfig struct {
	APIURL string `mapstructure:"api_url"`
}

type OpenAIConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type OpenClawConfig struct {
	GatewayURL string `mapstructure:"gateway_url"`
	AuthSecret string `mapstructure:"auth_secret"`
}

type RivetConfig struct {
	ServiceURL string `mapstructure:"service_url"`
}

type PostgresConfig struct {
	URL string `mapstructure:"url"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type RedpandaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

type HasuraConfig struct {
	URL         string `mapstructure:"url"`
	AdminSecret string `mapstructure:"admin_secret"`
}

type AgentsConfig struct {
	Concurrency int `mapstructure:"concurrency"`
}

// Global config instance
var AppConfig Config

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", "3000")
	viper.SetDefault("redpanda.brokers", []string{"localhost:9092"})
	viper.SetDefault("redpanda.topic", "nexus_tasks")
	viper.SetDefault("redpanda.group_id", "nexus_worker_group")
	viper.SetDefault("agents.concurrency", 10)
	viper.SetDefault("postgres.url", "postgres://postgres:password@localhost:5432/chatwoot_dev?sslmode=disable")
	viper.SetDefault("rivet.service_url", "localhost:50051")
	viper.SetDefault("benthos.api_url", "http://wasm-processor:4195")
	viper.SetDefault("role", "monolith")

	defaultMCPServers := []map[string]interface{}{
		{
			"id":      "svelte-mcp",
			"name":    "Svelte MCP",
			"type":    "stdio",
			"command": "npx",
			"args":    []string{"-y", "svelte-mcp-server"},
		},
	}
	viper.SetDefault("mcp.servers", defaultMCPServers)

	viper.SetEnvPrefix("NEXUS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Dynamic Config: Watch for changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(&AppConfig); err != nil {
			log.Printf("Error reloading config: %v", err)
		} else {
			log.Printf("Config reloaded. New concurrency: %d", AppConfig.Agents.Concurrency)
		}
	})
	viper.WatchConfig()

	return &AppConfig, nil
}
