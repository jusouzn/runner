package invoker

import (
	"bytes"
	"net"
	"strconv"
	"strings"
	"testing"
)

// portaLivre devolve uma porta TCP que estava livre no momento da chamada.
func portaLivre(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("não foi possível obter porta livre: %v", err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestHTTP_OperacaoInvalida(t *testing.T) {
	_, err := HTTP("http://localhost:0", Params{Operacao: "exportar"})
	if err == nil {
		t.Fatal("esperava erro para operação inválida")
	}
	if !strings.Contains(err.Error(), "operação inválida") {
		t.Errorf("mensagem inesperada: %v", err)
	}
}

func TestHTTP_ConexaoRecusada(t *testing.T) {
	// Porta livre (sem servidor): a conexão deve ser recusada com mensagem clara.
	porta := portaLivre(t)
	_, err := HTTP(addr(porta), Params{Operacao: "assinar", Conteudo: "x", Pin: "1234"})
	if err == nil {
		t.Fatal("esperava erro de conexão recusada")
	}
	if !strings.Contains(err.Error(), "erro ao conectar ao servidor") {
		t.Errorf("mensagem inesperada: %v", err)
	}
}

func TestExtractError_ParseiaJsonDoJar(t *testing.T) {
	stderr := *bytes.NewBufferString(
		`{"status":"error","error_code":"ERR_INVALID_PARAM","message":"pin curto","details":["min 4"]}`)
	err := extractError(stderr)
	if err == nil || !strings.Contains(err.Error(), "pin curto") {
		t.Errorf("esperava mensagem 'pin curto', obtido: %v", err)
	}
	if !strings.Contains(err.Error(), "min 4") {
		t.Errorf("esperava detalhes na mensagem, obtido: %v", err)
	}
}

func TestExtractError_StderrVazio(t *testing.T) {
	err := extractError(bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "sem mensagem") {
		t.Errorf("esperava erro genérico para stderr vazio, obtido: %v", err)
	}
}

func addr(porta int) string {
	return "http://localhost:" + strconv.Itoa(porta)
}
