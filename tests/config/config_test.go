package config_test

import (
	"os"
	"testing"
	"io/ioutil"
	"encoding/json"
	"github.com/a-digi/coco-server/server"
)

func TestLoadConfig_Success(t *testing.T) {
	cfg := server.Config{
		Port: 1234,
		PidFile: "test.pid",
		DomainFiles: []string{"domain1.json"},
	}
	file, err := ioutil.TempFile("", "config_*.json")
	if err != nil {
		t.Fatalf("TempFile error: %v", err)
	}
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	if err := enc.Encode(cfg); err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	file.Close()
	loaded, err := server.LoadConfig(file.Name())
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if loaded.Port != cfg.Port || loaded.PidFile != cfg.PidFile {
		t.Errorf("Loaded config does not match")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := server.LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("Expected error for missing file")
	}
}
