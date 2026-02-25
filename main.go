package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/errpunk/monocloud-config/config"
)

func update(configPath string) error {
	log.Println("Fetching remote subscription...")
	remote, err := config.Fetch()
	if err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}
	log.Printf("Fetched: %d proxies, %d proxy-groups, %d rules",
		len(remote.Proxies), len(remote.ProxyGroups), len(remote.Rules))

	if err := config.Merge(configPath, remote); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}
	log.Printf("Merged into %s successfully", configPath)
	return nil
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	intervalStr := os.Getenv("UPDATE_INTERVAL")
	if intervalStr == "" {
		intervalStr = "1h"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("Invalid UPDATE_INTERVAL %q: %v", intervalStr, err)
	}

	log.Printf("Starting monocloud-config daemon (config=%s, interval=%s)", configPath, interval)

	// Run immediately on startup
	if err := update(configPath); err != nil {
		log.Printf("Initial update error: %v", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Graceful shutdown on SIGTERM / SIGINT (Docker stop sends SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-ticker.C:
			if err := update(configPath); err != nil {
				log.Printf("Update error: %v", err)
			}
		case sig := <-quit:
			log.Printf("Received signal %s, shutting down", sig)
			return
		}
	}
}
