package server

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/a-digi/coco-logger/logger"
)

// IsPortBusy prüft, ob der Port bereits belegt ist. Gibt ggf. die PID des Prozesses zurück (immer 0, da Go das nicht nativ unterstützt).
func IsPortBusy(port int) (bool, int, error) {
	addr := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// Port ist belegt
		return true, 0, nil
	}
	ln.Close()
	return false, 0, nil
}

// WaitForPortFreed wartet, bis der Port nicht mehr belegt ist (Timeout: 10s)
func WaitForPortFreed(port int, log logger.Logger) {
	for i := 0; i < 100; i++ {
		busy, _, _ := IsPortBusy(port)
		if !busy {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Warning("Port %d ist nach Timeout immer noch belegt!", port)
}

// killServerProcess sucht den Prozess, der den Port belegt, und killt ihn per SIGTERM (nur Unix).
func killServerProcess(port int, log logger.Logger) {
	// Finde PID mit lsof
	cmd := exec.Command("lsof", "-t", "-i", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		log.Warning("Konnte keinen Prozess für Port %d finden: %v", port, err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		pid, err := strconv.Atoi(strings.TrimSpace(line))

		if err != nil {
			log.Warning("Konnte PID nicht parsen: %v", err)
			continue
		}

		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			log.Warning("Konnte Prozess %d nicht beenden: %v", pid, err)
		} else {
			log.Warning("SIGTERM an Prozess %d (Port %d) gesendet.", pid, port)
		}
	}
}

// EnsurePortIsFree prüft, ob der Port frei ist. Falls nicht, wird der Prozess per Port-Suche gekillt und gewartet, bis der Port frei ist.
func EnsurePortIsFree(port int, log logger.Logger) {
	busy, _, _ := IsPortBusy(port)
	if busy {
		log.Warning("Port %d ist bereits belegt. Versuche, den Prozess zu beenden...", port)
		killServerProcess(port, log)
		WaitForPortFreed(port, log)
	}
}
