package main

import (
	"context"
	"crypto/tls"
	"ecowitt-proxy/local/config"
	"ecowitt-proxy/local/splitter"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// Get the directory of the executable
	exeDir := filepath.Dir("./")

	// Construct the full path to the config.yml file
	configPath := filepath.Join(exeDir, "config.yml")

	// Load the configuration
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Print the loaded configuration
	fmt.Printf("Config: %+v\n", cfg)

	// Create the HTTP client
	client := buildClient(cfg)

	// Initialize the splitter
	s := splitter.Splitter{
		Config: cfg,
		Client: client,
	}

	// Set up the HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Server.Path, s.HandleRequest)

	const listenPort = 8123
	listenAddr := fmt.Sprintf(":%d", listenPort)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Starting proxy server on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
		}
	}()

	// Wait for interrupt signal
	<-stop

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	log.Println("Shutting down the proxy server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Proxy server forced to shutdown: %v", err)
	}

	log.Println("Proxy server stopped")
}

// buildClient creates an HTTP client with a timeout and custom TLS settings
func buildClient(cfg config.Config) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// Custom TLS settings (not recommended for production if InsecureSkipVerify is true)
				InsecureSkipVerify: cfg.Server.TlsInsecureSkipVerify,
			},
		},
	}
}
