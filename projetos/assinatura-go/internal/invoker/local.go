package invoker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Params são os parâmetros para invocar o assinador.jar.
type Params struct {
	Operacao   string // "assinar" ou "validar"
	Conteudo   string
	Pin        string
	Algoritmo  string // opcional; padrão: SHA256withRSA
	Assinatura string // obrigatório somente em "validar"
}

// Local invoca o assinador.jar diretamente via `java -jar` (cold start).
// Localiza o java e o jar automaticamente antes de executar.
func Local(p Params) (string, error) {
	javaPath, err := findJava()
	if err != nil {
		return "", err
	}
	jarPath, err := findJar()
	if err != nil {
		return "", err
	}

	args := buildArgs(jarPath, p)
	cmd := exec.Command(javaPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", extractError(stderr)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// buildArgs monta a lista de argumentos para `java`.
func buildArgs(jarPath string, p Params) []string {
	args := []string{
		"-jar", jarPath,
		"--modo=local",
		"--operacao=" + p.Operacao,
		"--pin=" + p.Pin,
		"--conteudo=" + p.Conteudo,
	}
	if p.Algoritmo != "" {
		args = append(args, "--algoritmo="+p.Algoritmo)
	}
	if p.Assinatura != "" {
		args = append(args, "--assinatura="+p.Assinatura)
	}
	return args
}

// extractError interpreta o stderr do assinador.jar.
// Tenta extrair a mensagem do JSON de erro; se não conseguir, retorna o texto bruto.
func extractError(stderr bytes.Buffer) error {
	if stderr.Len() == 0 {
		return errors.New("assinador.jar encerrou com erro sem mensagem")
	}
	var errResp struct {
		Message string   `json:"message"`
		Details []string `json:"details"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &errResp); err == nil && errResp.Message != "" {
		if len(errResp.Details) > 0 {
			return fmt.Errorf("%s: %s", errResp.Message, strings.Join(errResp.Details, "; "))
		}
		return errors.New(errResp.Message)
	}
	return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
}

// findJava localiza o binário java: primeiro em ~/.hubsaude/jdk, depois no PATH.
func findJava() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		managed := filepath.Join(home, ".hubsaude", "jdk", "bin", "java")
		if _, err := os.Stat(managed); err == nil {
			return managed, nil
		}
	}
	if path, err := exec.LookPath("java"); err == nil {
		return path, nil
	}
	return "", errors.New(
		"java não encontrado — instale o Java 21 ou execute `assinatura setup` para provisionamento automático",
	)
}

// findJar localiza o assinador.jar: ~/.hubsaude/, ao lado do executável ou diretório atual.
func findJar() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		managed := filepath.Join(home, ".hubsaude", "assinador.jar")
		if _, err := os.Stat(managed); err == nil {
			return managed, nil
		}
	}
	if exe, err := os.Executable(); err == nil {
		next := filepath.Join(filepath.Dir(exe), "assinador.jar")
		if _, err := os.Stat(next); err == nil {
			return next, nil
		}
	}
	// diretório atual (útil durante o desenvolvimento)
	if _, err := os.Stat("assinador.jar"); err == nil {
		return "assinador.jar", nil
	}
	return "", errors.New(
		"assinador.jar não encontrado — coloque-o em ~/.hubsaude/ ou no mesmo diretório do executável",
	)
}
