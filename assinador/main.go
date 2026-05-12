// main.go - CLI `assinatura` do Sistema Runner.
//
// Suporta dois modos de invocacao do assinador.jar:
//   - modo HTTP (warm start, padrao):     CLI -> http://localhost:porta/sign
//   - modo local (cold start, --local):   CLI -> java -jar assinador.jar
//
// O contrato de mensagens esta descrito em Design.md e docs/INTEGRACAO.md.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION      = "0.2.0"
	PORTA_PADRAO = 8080
	JAR_PADRAO   = "assinador.jar"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--version", "-v":
		fmt.Println(VERSION)
	case "--help", "-h":
		printUsage()
	case "assinar":
		cmdOperacao("assinar", os.Args[2:])
	case "validar":
		cmdOperacao("validar", os.Args[2:])
	case "servidor":
		cmdServidor(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Comando '%s' nao reconhecido.\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`assinatura - CLI do Sistema Runner

Uso:
  assinatura assinar [--local] --pin=NNNN [--porta=N] [--jar=PATH]
  assinatura validar [--local] --pin=NNNN [--assinatura=...] [--porta=N] [--jar=PATH]
  assinatura servidor iniciar [--porta=N] [--jar=PATH]
  assinatura servidor parar   [--porta=N]
  assinatura servidor status  [--porta=N]
  assinatura --version | --help

Modos:
  Sem --local, o CLI utiliza o assinador.jar em modo servidor (HTTP).
    - reutiliza instancia ja em execucao na porta;
    - inicia automaticamente se nao houver instancia ativa.
  Com --local, executa diretamente "java -jar assinador.jar ..." (cold start).`)
}

// ---------- Flags ----------

type Flags struct {
	Local      bool
	Porta      int
	Jar        string
	Pin        string
	Assinatura string
}

func parseFlags(args []string) *Flags {
	f := &Flags{Porta: PORTA_PADRAO, Jar: JAR_PADRAO}
	for _, a := range args {
		switch {
		case a == "--local":
			f.Local = true
		case strings.HasPrefix(a, "--porta="):
			if n, err := strconv.Atoi(strings.TrimPrefix(a, "--porta=")); err == nil && n > 0 {
				f.Porta = n
			}
		case strings.HasPrefix(a, "--jar="):
			f.Jar = strings.TrimPrefix(a, "--jar=")
		case strings.HasPrefix(a, "--pin="):
			f.Pin = strings.TrimPrefix(a, "--pin=")
		case strings.HasPrefix(a, "--assinatura="):
			f.Assinatura = strings.TrimPrefix(a, "--assinatura=")
		}
	}
	return f
}

// ---------- Operacoes (assinar / validar) ----------

func cmdOperacao(op string, args []string) {
	f := parseFlags(args)
	if f.Pin == "" {
		printErroCLI("ERR_INVALID_PARAM", "Parametro obrigatorio ausente: --pin")
		os.Exit(1)
	}
	if f.Local {
		executarLocal(op, f)
	} else {
		executarHTTP(op, f)
	}
}

func executarLocal(op string, f *Flags) {
	fmt.Fprintf(os.Stderr, "[CLI] Executando assinador.jar em modo local (%s)\n", op)
	args := []string{"-jar", f.Jar, "--modo=local", "--operacao=" + op, "--pin=" + f.Pin}
	if f.Assinatura != "" {
		args = append(args, "--assinatura="+f.Assinatura)
	}
	cmd := exec.Command("java", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}

func executarHTTP(op string, f *Flags) {
	if !servidorAtivo(f.Porta) {
		if err := iniciarServidor(f); err != nil {
			printErroCLI("ERR_SERVIDOR_INDISPONIVEL", err.Error())
			os.Exit(1)
		}
	}

	endpoint := map[string]string{"assinar": "sign", "validar": "validate"}[op]
	body, _ := json.Marshal(map[string]string{
		"pin":        f.Pin,
		"assinatura": f.Assinatura,
	})
	url := fmt.Sprintf("http://localhost:%d/%s", f.Porta, endpoint)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		printErroCLI("ERR_HTTP", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}

// ---------- Subcomando "servidor" ----------

func cmdServidor(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Subcomando necessario: iniciar | parar | status")
		os.Exit(1)
	}
	sub := args[0]
	f := parseFlags(args[1:])

	switch sub {
	case "iniciar":
		if servidorAtivo(f.Porta) {
			fmt.Printf("Servidor ja em execucao na porta %d.\n", f.Porta)
			return
		}
		if err := iniciarServidor(f); err != nil {
			printErroCLI("ERR_SERVIDOR_INDISPONIVEL", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Servidor iniciado na porta %d.\n", f.Porta)

	case "parar":
		cli := http.Client{Timeout: 3 * time.Second}
		url := fmt.Sprintf("http://localhost:%d/shutdown", f.Porta)
		resp, err := cli.Post(url, "application/json", nil)
		if err != nil {
			printErroCLI("ERR_SERVIDOR_INDISPONIVEL", err.Error())
			os.Exit(1)
		}
		resp.Body.Close()
		limparEstado()
		fmt.Printf("Servidor na porta %d parado.\n", f.Porta)

	case "status":
		if servidorAtivo(f.Porta) {
			if e, ok := lerEstado(); ok && e.Porta == f.Porta {
				fmt.Printf("Servidor ATIVO na porta %d (PID %d).\n", e.Porta, e.PID)
			} else {
				fmt.Printf("Servidor ATIVO na porta %d.\n", f.Porta)
			}
		} else {
			fmt.Printf("Servidor INATIVO (porta %d).\n", f.Porta)
		}

	default:
		fmt.Fprintf(os.Stderr, "Subcomando '%s' nao reconhecido.\n", sub)
		os.Exit(1)
	}
}

// ---------- Helpers de servidor ----------

func servidorAtivo(porta int) bool {
	cli := http.Client{Timeout: 800 * time.Millisecond}
	resp, err := cli.Get(fmt.Sprintf("http://localhost:%d/health", porta))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func portaOcupada(porta int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", porta))
	if err != nil {
		return true
	}
	_ = ln.Close()
	return false
}

func iniciarServidor(f *Flags) error {
	if portaOcupada(f.Porta) {
		// Pode ser o proprio assinador subindo; sanity-check via /health.
		if servidorAtivo(f.Porta) {
			return nil
		}
		return fmt.Errorf("porta %d ocupada por outro processo", f.Porta)
	}

	fmt.Fprintf(os.Stderr, "[CLI] Iniciando assinador.jar em modo servidor (porta %d)...\n", f.Porta)
	cmd := exec.Command("java", "-jar", f.Jar,
		"--modo=server",
		"--porta="+strconv.Itoa(f.Porta),
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	salvarEstado(cmd.Process.Pid, f.Porta)

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if servidorAtivo(f.Porta) {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("servidor nao respondeu /health no tempo esperado")
}

// ---------- Estado local (~/.hubsaude/assinador.json) ----------

type estado struct {
	PID   int `json:"pid"`
	Porta int `json:"porta"`
}

func arquivoEstado() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".hubsaude")
	_ = os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "assinador.json")
}

func salvarEstado(pid, porta int) {
	f, err := os.Create(arquivoEstado())
	if err != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(estado{PID: pid, Porta: porta})
}

func limparEstado() { _ = os.Remove(arquivoEstado()) }

func lerEstado() (estado, bool) {
	var e estado
	f, err := os.Open(arquivoEstado())
	if err != nil {
		return e, false
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return e, false
	}
	return e, true
}

// ---------- Saida estruturada de erro do CLI ----------

func printErroCLI(code, msg string) {
	e := map[string]any{
		"status":     "error",
		"error_code": code,
		"message":    msg,
		"details":    []string{},
	}
	b, _ := json.MarshalIndent(e, "", "  ")
	fmt.Fprintln(os.Stderr, string(b))
}
