package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-digi/coco-logger/logger"
	"github.com/a-digi/coco-server/server/routing"
	"github.com/a-digi/coco-server/server/security"
)

func StartServer(configPath string, log logger.Logger) (*http.Server, *Config, error) {
	log.Info("Starting StartServer()...")
	config, err := ensureConfig(configPath)

	if err != nil {
		log.Error("Failed to load or create config: %v", err)
		return nil, nil, err
	}

	log.Info("Loaded config: %+v", config)

	EnsurePortIsFree(config.Port, log)

	if err := WritePID(config.PidFile); err != nil {
		log.Error("Failed to write PID file: %v", err)
		return nil, nil, err
	}
	log.Info("Wrote PID file to %s", config.PidFile)

	addr := fmt.Sprintf(":%d", config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: routing.GlobalRouteBuilder.Build(log),
	}

	log.Info("Starting server on %s", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed: %v", err)
		}
	}()

	return server, config, nil
}

// GracefulShutdown wartet immer selbst auf SIGTERM/SIGINT und loggt das Signal.
func GracefulShutdown(server *http.Server, pidFile string, log logger.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	sig := <-ch
	log.Info("Received signal %s, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			log.Error("Graceful shutdown failed: %v", err)
		} else {
			log.Info("Server gracefully stopped.")
		}
	}

	if err := RemovePID(pidFile); err != nil {
		log.Warning("Could not remove PID file: %v", err)
	} else {
		log.Info("Removed PID file %s", pidFile)
	}
}

func ShutdownServer(configPath string, log logger.Logger) error {
	config, err := LoadConfig(configPath)

	if err != nil {
		log.Error("Could not load config: %v", err)
		return err
	}

	pidFile := config.PidFile

	if pidFile == "" {
		pidFile = "./server.pid"
	}

	pid, err := ReadPID(pidFile)

	if err != nil {
		log.Error("Could not read PID file: %v", err)
		return err
	}

	err = syscall.Kill(pid, syscall.SIGTERM)

	if err != nil {
		log.Error("Failed to send SIGTERM: %v", err)
		return err
	}

	log.Info("Shutdown signal sent.")
	return nil
}

func SendSIGTERM(pid int) error {
	err := syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}

// Add a setter for the security layer
func SetSecurityLayer(layer security.SecurityLayer) {
	routing.GlobalRouteBuilder.SetSecurityLayer(layer)
}
