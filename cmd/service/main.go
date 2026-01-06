package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Get the target URL from environment variable, default to internal docker name
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		fmt.Println("TARGET_URL is not set, defaulting to http://wordpress")
		targetURL = "http://wordpress"
	}

	fmt.Printf("Checking connection to: %s\n", targetURL)

	resp, err := http.Get(targetURL) // #nosec G107
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", targetURL, err)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	fmt.Printf("Success! Connected to %s. Status: %s\n", targetURL, resp.Status)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Shutting down...")
}
