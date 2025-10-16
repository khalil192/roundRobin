package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type AppConfig struct {
	BackendServerURLs []string
}

var cfg *AppConfig

func Config(fileName string) *AppConfig {
	if cfg == nil {
		loadConfig(fileName)
	}

	return cfg
}

func loadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to open config file: %w", err))
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			servers = append(servers, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error reading from file: %w", err))
	}

	cfg = &AppConfig{BackendServerURLs: servers}
	fmt.Println("Loaded backend servers:", cfg.BackendServerURLs)
}
