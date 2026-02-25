package main

import (
	"fmt"
	"os"

	"github.com/errpunk/monocloud-config/config"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	fmt.Println("Fetching remote subscription...")
	remote, err := config.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Fetched: %d proxies, %d proxy-groups, %d rules\n",
		len(remote.Proxies), len(remote.ProxyGroups), len(remote.Rules))

	fmt.Printf("Merging into %s...\n", configPath)
	if err := config.Merge(configPath, remote); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}
