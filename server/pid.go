package server

import (
	"fmt"
	"os"
)

// WritePID writes the current process PID to the specified file
func WritePID(pidFile string) error {
	pid := os.Getpid()
	f, err := os.Create(pidFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%d", pid))
	return err
}

// ReadPID reads the PID from the specified file
func ReadPID(pidFile string) (int, error) {
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	var pid int
	_, err = fmt.Sscanf(string(pidData), "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// RemovePID removes the PID file
func RemovePID(pidFile string) error {
	return os.Remove(pidFile)
}
