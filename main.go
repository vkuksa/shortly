package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/vkuksa/shortly/assets"
	shortly "github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/http"
	"github.com/vkuksa/shortly/internal/shortener"
	"github.com/vkuksa/shortly/pkg/storage"
	"github.com/vkuksa/shortly/pkg/storage/inmem"
)

const (
	// DefaultConfigPath is the default path to the application configuration.
	DefaultConfigPath = "shortly.conf"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Parse command line flags
	configPath, err := parseFlags(os.Args[1:])
	if err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		log.Fatalf(err.Error())
	}

	// Load config
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Instantiate a new type to represent application.
	m := NewMain(config)

	// Execute program.
	if err := m.Run(); err != nil {
		m.Close()
		log.Fatalf(err.Error())
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		log.Fatalf(err.Error())
	}
}

func parseFlags(args []string) (configPath string, err error) {
	fs := flag.NewFlagSet("shortly", flag.ContinueOnError)
	fs.StringVar(&configPath, "config", DefaultConfigPath, "config path")
	err = fs.Parse(args)
	return
}

// Main represents the program.
type Main struct {
	Config *Config

	Storage storage.Storage[shortly.Link]
	Service shortly.LinkService

	HTTPServer *http.Server
}

func NewMain(c *Config) *Main {
	storage := NewStorage[shortly.Link](c.DS.Kind)
	service := shortener.NewService(storage)
	server := http.NewServer()

	server.LinkService = service
	server.Addr = c.HTTP.Addr
	server.Scheme = c.HTTP.Scheme
	server.Domain = c.HTTP.Domain
	server.Assets = assets.All

	return &Main{
		Config: c,

		Service:    service,
		Storage:    storage,
		HTTPServer: server,
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.Storage != nil {
		if err := m.Storage.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Run executes the program
func (m *Main) Run() (err error) {
	return m.HTTPServer.Open()
}

type Config struct {
	DS struct {
		Kind string `toml:"kind"`
		Name string `toml:"name"`
	} `toml:"ds"`

	HTTP struct {
		Addr   string `toml:"addr"`
		Scheme string `toml:"scheme"`
		Domain string `toml:"domain"`
	} `toml:"http"`
}

func NewConfig(filepath string) (*Config, error) {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", filepath)
		}

		return nil, fmt.Errorf("newconfig: %w", err)
	}

	config := &Config{}
	if err := toml.Unmarshal(buf, config); err != nil {
		return nil, fmt.Errorf("newconfig: %w", err)
	}

	return config, nil
}

func NewStorage[V any](kind string) storage.Storage[V] {
	switch kind {
	case "inmem":
		return inmem.NewStorage[V]()
	default:
		log.Fatalf("Link storage %s is not supported", kind)
		return nil
	}
}
