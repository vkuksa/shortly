package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/vkuksa/shortly"
	"github.com/vkuksa/shortly/internal/business"
	"github.com/vkuksa/shortly/internal/http"
	"github.com/vkuksa/shortly/internal/storage/inmem"
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
	configPath, err := parseFlags(ctx, os.Args[1:])
	if err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Load config
	config, err := NewConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Instantiate a new type to represent application.
	m := NewMain(config)

	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseFlags(ctx context.Context, args []string) (configPath string, err error) {
	fs := flag.NewFlagSet("shortly", flag.ContinueOnError)
	fs.StringVar(&configPath, "config", DefaultConfigPath, "config path")
	err = fs.Parse(args)
	return
}

// Main represents the program.
type Main struct {
	Config *Config

	LinkStorage shortly.LinkStorage
	LinkService shortly.LinkService

	HTTPServer *http.Server
}

func NewMain(c *Config) *Main {
	storage := CreateLinkStorage(c.DS.Kind)
	service := business.NewLinkService(storage)
	server := http.NewServer()

	server.LinkService = service
	server.Addr = c.HTTP.Addr
	server.Scheme = c.HTTP.Scheme
	server.Domain = c.HTTP.Domain

	return &Main{
		Config: c,

		LinkService: service,
		LinkStorage: storage,
		HTTPServer:  server,
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.LinkStorage != nil {
		if err := m.LinkStorage.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Run executes the program
func (m *Main) Run(ctx context.Context) (err error) {
	if err := m.LinkStorage.Open(ctx); err != nil {
		return fmt.Errorf("cannot open data source: %w", err)
	}
	return m.HTTPServer.Open()
}

type Config struct { //TODO: define appropriate config?
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
		} else {
			return nil, fmt.Errorf("newconfig: %w", err)
		}
	}

	config := &Config{}
	if err := toml.Unmarshal(buf, config); err != nil {
		return nil, fmt.Errorf("newconfig: %w", err)
	}

	return config, nil
}

func CreateLinkStorage(kind string) shortly.LinkStorage {
	switch kind {
	case "inmem":
		return inmem.NewLinkStorage()
	default:
		panic(fmt.Sprintf("Link storage %s is not supported", kind))
	}
}
