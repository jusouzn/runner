package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/jusouzn/assinatura/internal/jdk"
)

const DefaultPort = 8080

type state struct {
	PID  int `json:"pid"`
	Port int `json:"port"`
}

// EnsureRunning verifies that assinador server is up on port.
// If already running, reuses it. If not, starts it and waits until ready.
func EnsureRunning(port int) error {
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if isAlive(baseURL) {
		return nil
	}
	return startServer(port)
}

// Stop sends SIGINT to the assinador server and removes the state file.
func Stop(port int) error {
	sf, err := stateFilePath()
	if err != nil {
		return err
	}
	s, err := loadState(sf)
	if err != nil || s.PID == 0 {
		return fmt.Errorf("servidor não está em execução (nenhum estado encontrado em %s)", sf)
	}
	if err := killPID(s.PID); err != nil {
		return fmt.Errorf("erro ao encerrar o servidor (PID %d): %w", s.PID, err)
	}
	os.Remove(sf)
	return nil
}

func isAlive(baseURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/health", nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func startServer(port int) error {
	javaPath, err := jdk.FindJava()
	if err != nil {
		return err
	}
	jarPath, err := jdk.FindJar()
	if err != nil {
		return err
	}

	cmd := exec.Command(javaPath, "-jar", jarPath,
		"--modo=servidor",
		"--porta="+strconv.Itoa(port),
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar assinador.jar: %w", err)
	}

	if sf, err := stateFilePath(); err == nil {
		os.MkdirAll(filepath.Dir(sf), 0755)
		saveState(sf, state{PID: cmd.Process.Pid, Port: port}) //nolint:errcheck
	}

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if err := waitReady(baseURL, 20*time.Second); err != nil {
		return fmt.Errorf("assinador.jar iniciado (PID %d) mas não respondeu: %w", cmd.Process.Pid, err)
	}
	return nil
}

func waitReady(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isAlive(baseURL) {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("servidor não respondeu em %s", timeout)
}

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hubsaude", "assinador-server.json"), nil
}

func loadState(path string) (state, error) {
	var s state
	data, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(data, &s)
}

func saveState(path string, s state) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func killPID(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return p.Kill()
	}
	return p.Signal(os.Interrupt)
}
