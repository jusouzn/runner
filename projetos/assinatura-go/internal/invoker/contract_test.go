//go:build contract

// Testes de contrato CLI ↔ assinador.jar: exercem o JAR real via subprocess
// (modo local) e via HTTP (modo servidor). Rodam apenas com a tag `contract`
// e exigem a variável de ambiente ASSINADOR_JAR apontando para o JAR e um
// java >= 21 no PATH.
//
//   go test -tags=contract ./...
package invoker

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jusouzn/assinatura/internal/jdk"
	"github.com/jusouzn/assinatura/internal/server"
)

// prepararAmbiente isola HOME e instala o JAR informado em ASSINADOR_JAR no
// caminho gerenciado (~/.hubsaude/assinador.jar). Faz skip se faltar JAR ou java.
func prepararAmbiente(t *testing.T) {
	t.Helper()
	origem := os.Getenv("ASSINADOR_JAR")
	if origem == "" {
		t.Skip("ASSINADOR_JAR não definido — pule os testes de contrato")
	}
	dados, err := os.ReadFile(origem)
	if err != nil {
		t.Skipf("não foi possível ler ASSINADOR_JAR (%s): %v", origem, err)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	destDir := filepath.Join(home, ".hubsaude")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "assinador.jar"), dados, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := jdk.FindJava(); err != nil {
		t.Skipf("java >= %d indisponível: %v", jdk.MinJavaVersion, err)
	}
}

func portaLivreContrato(t *testing.T) int {
	t.Helper()
	return portaLivre(t) // reutiliza helper de invoker_negative_test.go
}

func TestContrato_LocalAssinarEValidar(t *testing.T) {
	prepararAmbiente(t)

	out, err := Local(Params{Operacao: "assinar", Conteudo: "olá, mundo com acentos", Pin: "1234"})
	if err != nil {
		t.Fatalf("assinar falhou: %v", err)
	}
	if !strings.Contains(out, "\"status\": \"success\"") {
		t.Errorf("saída inesperada de assinar: %s", out)
	}

	out, err = Local(Params{Operacao: "validar", Conteudo: "conteudo", Assinatura: "abc", Pin: "1234"})
	if err != nil {
		t.Fatalf("validar falhou: %v", err)
	}
	if !strings.Contains(out, "\"valid\": true") {
		t.Errorf("saída inesperada de validar: %s", out)
	}
}

func TestContrato_LocalPropagaErroDeParametro(t *testing.T) {
	prepararAmbiente(t)

	// PIN inválido: o JAR deve sair com erro e a CLI propagar a mensagem,
	// sem produzir resultado em stdout.
	out, err := Local(Params{Operacao: "assinar", Conteudo: "x", Pin: "12"})
	if err == nil {
		t.Fatal("esperava erro para PIN inválido")
	}
	if out != "" {
		t.Errorf("stdout deveria ser vazio em caso de erro, obtido: %q", out)
	}
	if !strings.Contains(err.Error(), "pin") && !strings.Contains(err.Error(), "PIN") {
		t.Errorf("mensagem deveria mencionar o PIN: %v", err)
	}
}

func TestContrato_ServidorHTTPeIdempotencia(t *testing.T) {
	prepararAmbiente(t)

	porta := portaLivreContrato(t)
	if err := server.EnsureRunning(porta); err != nil {
		t.Fatalf("EnsureRunning falhou: %v", err)
	}
	t.Cleanup(func() { _ = server.Stop(porta) })

	// Idempotência: segunda chamada reaproveita a instância viva (health check real).
	inicio := time.Now()
	if err := server.EnsureRunning(porta); err != nil {
		t.Fatalf("segunda EnsureRunning falhou: %v", err)
	}
	if elapsed := time.Since(inicio); elapsed > 3*time.Second {
		t.Errorf("segunda chamada demorou demais (%s) — não reaproveitou a instância", elapsed)
	}

	baseURL := addr(porta)
	out, err := HTTP(baseURL, Params{Operacao: "assinar", Conteudo: "via http", Pin: "1234"})
	if err != nil {
		t.Fatalf("requisição HTTP de assinatura falhou: %v", err)
	}
	if !strings.Contains(out, "\"status\": \"success\"") {
		t.Errorf("resposta HTTP inesperada: %s", out)
	}
}

func TestContrato_PortaOcupadaFalhaClara(t *testing.T) {
	prepararAmbiente(t)

	// Ocupa uma porta com um listener que não fala HTTP de health.
	porta := portaLivre(t)
	ln, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(porta))
	if err != nil {
		t.Fatalf("não foi possível ocupar a porta: %v", err)
	}
	defer ln.Close()

	err = server.EnsureRunning(porta)
	if err == nil {
		_ = server.Stop(porta)
		t.Fatal("esperava falha clara ao subir com a porta ocupada")
	}
	if !strings.Contains(err.Error(), "ocupada") && !strings.Contains(err.Error(), "não ficou pronto") {
		t.Errorf("mensagem deveria indicar porta ocupada/não pronto: %v", err)
	}
}
