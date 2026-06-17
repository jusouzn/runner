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

// StopOrSchedule encerra o servidor imediatamente (aposMinutos==0) ou agenda o
// encerramento por inatividade via POST /shutdown?apos=N.
// Tenta HTTP primeiro; cai em kill por PID apenas para encerramento imediato.
func StopOrSchedule(port, aposMinutos int) error {
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	shutdownURL := baseURL + "/shutdown"
	if aposMinutos > 0 {
		shutdownURL = fmt.Sprintf("%s?apos=%d", shutdownURL, aposMinutos)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, shutdownURL, nil)
	if err == nil {
		if resp, doErr := http.DefaultClient.Do(req); doErr == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				if aposMinutos == 0 {
					if sf, sfErr := stateFilePath(); sfErr == nil {
						os.Remove(sf)
					}
				}
				return nil
			}
		}
	}

	if aposMinutos > 0 {
		return fmt.Errorf("servidor não está acessível via HTTP na porta %d; encerramento agendado não é possível", port)
	}
	return Stop(port)
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

	// Observa o término do processo para falhar rápido (e com mensagem clara)
	// quando o JAR não consegue subir — tipicamente porta já ocupada — em vez
	// de esperar o timeout inteiro de readiness.
	exited := make(chan struct{})
	go func() {
		cmd.Wait() //nolint:errcheck
		close(exited)
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	if err := waitReady(baseURL, exited, 20*time.Second); err != nil {
		// Cleanup best-effort: o state file pode estar apontando para um PID
		// morto (ou o processo pode seguir rodando após timeout). Removemos o
		// arquivo e tentamos encerrar o processo para evitar estado
		// inconsistente em chamadas futuras de Stop/EnsureRunning.
		pid := cmd.Process.Pid
		killPID(pid) //nolint:errcheck
		if sf, sfErr := stateFilePath(); sfErr == nil {
			os.Remove(sf)
		}
		return fmt.Errorf("assinador.jar (PID %d) não ficou pronto na porta %d: %w", pid, port, err)
	}
	return nil
}

// waitReady aguarda o servidor responder ao health check. Retorna erro
// imediatamente se o processo terminar antes (porta ocupada, JVM falhou etc.).
func waitReady(baseURL string, exited <-chan struct{}, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isAlive(baseURL) {
			return nil
		}
		select {
		case <-exited:
			return fmt.Errorf("o processo encerrou antes de ficar pronto (porta %s pode estar ocupada)",
				baseURL[len("http://localhost:"):])
		case <-time.After(300 * time.Millisecond):
		}
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
