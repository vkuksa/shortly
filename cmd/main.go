package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	shortly "github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/http"
	"github.com/vkuksa/shortly/internal/shortener"
	"github.com/vkuksa/shortly/pkg/storage"
	"github.com/vkuksa/shortly/pkg/storage/bbolt"
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
		log.Fatal(err.Error())
	}

	// Load config
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiate a new type to represent application.
	m, err := NewMain(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Execute program.
	if err := m.Run(); err != nil {
		log.Print(err.Error())
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		log.Print(err.Error())
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

func NewMain(c *Config) (*Main, error) {
	storage, err := NewStorage[shortly.Link](c.DB)
	if err != nil {
		return nil, fmt.Errorf("NewMain: %w", err)
	}

	service := shortener.NewService(storage)
	server := http.NewServer(c.HTTP)

	server.LinkService = service

	return &Main{
		Config: c,

		Service:    service,
		Storage:    storage,
		HTTPServer: server,
	}, nil
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
	HTTP http.Config `toml:"http"`

	DB DBConfig `toml:"db"`
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

type DBConfig struct {
	Kind       string `toml:"kind"`
	Connection struct {
		Host     string `toml:"host"`
		Port     int64  `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	}
	BBolt struct {
		File   string `toml:"file"`
		Bucket string `toml:"bucket"`
	} `toml:"bbolt"`
}

func NewStorage[V any](c DBConfig) (storage.Storage[V], error) {
	switch c.Kind {
	case "inmem":
		return inmem.NewStorage[V](), nil
	case "bbolt":
		stor, err := bbolt.NewStorage[V](c.BBolt.File, c.BBolt.Bucket)
		if err != nil {
			return nil, fmt.Errorf("NewStorage: %w", err)
		}

		return stor, nil
	default:
		return nil, fmt.Errorf("NewStorage: Storage %s is not supported", c.Kind)
	}
}
