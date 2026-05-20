package process

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const DefaultPort = 8443

// insecureClient aceita certificados auto-assinados do simulador.
var insecureClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
	},
}

type State struct {
	PID  int `json:"pid"`
	Port int `json:"port"`
}

// IsPortFree verifica se a porta TCP está disponível para uso.
func IsPortFree(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// IsRunning verifica se o simulador está respondendo via TCP na porta indicada.
func IsRunning(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Start inicia o simulador.jar em background e salva o PID.
// Aguarda até 30s para o simulador aceitar conexões TCP.
func Start(javaPath, jarPath string, port int) error {
	cmd := exec.Command(javaPath, "-jar", jarPath)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar simulador.jar: %w", err)
	}

	if sf, err := stateFilePath(); err == nil {
		os.MkdirAll(filepath.Dir(sf), 0755)
		saveState(sf, State{PID: cmd.Process.Pid, Port: port}) //nolint:errcheck
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if IsRunning(port) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("simulador iniciado (PID %d) mas não respondeu em 30s", cmd.Process.Pid)
}

// Stop encerra o simulador: tenta o endpoint /shutdown; se falhar, mata o processo pelo PID.
func Stop(port int) error {
	// Tentativa de graceful shutdown via endpoint HTTP/HTTPS
	for _, scheme := range []string{"https", "http"} {
		url := fmt.Sprintf("%s://localhost:%d/shutdown", scheme, port)
		req, err := http.NewRequestWithContext(
			context.Background(), http.MethodPost, url, nil)
		if err != nil {
			continue
		}
		resp, err := insecureClient.Do(req)
		if err == nil {
			resp.Body.Close()
			removeState()
			return nil
		}
	}

	// Fallback por PID
	sf, err := stateFilePath()
	if err != nil {
		return fmt.Errorf("não foi possível localizar o arquivo de estado")
	}
	s, err := loadState(sf)
	if err != nil || s.PID == 0 {
		return fmt.Errorf("simulador não está em execução (nenhum estado encontrado em %s)", sf)
	}
	if err := killPID(s.PID); err != nil {
		return fmt.Errorf("erro ao encerrar simulador (PID %d): %w", s.PID, err)
	}
	os.Remove(sf)
	return nil
}

// Status consulta o endpoint /api/info e retorna o JSON formatado.
func Status(port int) (string, error) {
	for _, scheme := range []string{"https", "http"} {
		url := fmt.Sprintf("%s://localhost:%d/api/info", scheme, port)
		req, err := http.NewRequestWithContext(
			context.Background(), http.MethodGet, url, nil)
		if err != nil {
			continue
		}
		resp, err := insecureClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var info map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return fmt.Sprintf("status HTTP %d (corpo não é JSON)", resp.StatusCode), nil
		}
		data, _ := json.MarshalIndent(info, "", "  ")
		return string(data), nil
	}
	return "", fmt.Errorf("simulador não está respondendo na porta %d", port)
}

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hubsaude", "simulador-server.json"), nil
}

func loadState(path string) (State, error) {
	var s State
	data, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(data, &s)
}

func saveState(path string, s State) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func removeState() {
	if sf, err := stateFilePath(); err == nil {
		os.Remove(sf)
	}
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
