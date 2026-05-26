package invoker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jusouzn/assinatura/internal/jdk"
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
func Local(p Params) (string, error) {
	javaPath, err := jdk.FindJava()
	if err != nil {
		return "", err
	}
	jarPath, err := jdk.FindJar()
	if err != nil {
		return "", err
	}

	args := buildArgs(jarPath, p)

	if Verbose {
		logVerbose("executando: %s %s", javaPath, strings.Join(args, " "))
	}

	cmd := exec.Command(javaPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", extractError(stderr)
	}
	return strings.TrimSpace(stdout.String()), nil
}

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

// extractError interprets the stderr of assinador.jar.
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
