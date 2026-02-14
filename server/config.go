package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const DefaultPort = 2026

type Config struct {
	Port        int      `json:"port"`
	PidFile     string   `json:"pid_file"`
	DomainFiles []string `json:"domain_files,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &config, nil
}

func ensureConfig(path string) (*Config, error) {
	// Early return if config exists
	info, statErr := os.Stat(path)
	if statErr == nil && info != nil {
		return LoadConfig(path)
	}
	if statErr != nil && !os.IsNotExist(statErr) {
		return nil, fmt.Errorf("could not stat config file: %w", statErr)
	}

	log.Printf("[WARN] Config file %s not found. Creating default config...", path)
	// Use the current working directory for the default PID file
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get working directory: %w", err)
	}
	defaultPidFile := filepath.Join(cwd, "server.pid")
	defaultConfig := Config{
		Port:        DefaultPort,
		DomainFiles: nil,
		PidFile:     defaultPidFile,
	}
	data, marshalErr := json.MarshalIndent(defaultConfig, "", "  ")
	if marshalErr != nil {
		return nil, fmt.Errorf("could not marshal default config: %w", marshalErr)
	}
	if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
		return nil, fmt.Errorf("could not write default config: %w", writeErr)
	}
	return LoadConfig(path)
}
