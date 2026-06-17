package jdk

import (
	"strings"
	"testing"
)

// isolarAmbiente aponta HOME e PATH para um diretório vazio temporário,
// garantindo que nem o JDK gerenciado nem um java no PATH sejam encontrados.
func isolarAmbiente(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir) // HOME no Windows
	t.Setenv("PATH", dir)
}

func TestFindJar_Ausente(t *testing.T) {
	isolarAmbiente(t)
	_, err := FindJar()
	if err == nil {
		t.Fatal("esperava erro quando o assinador.jar não existe")
	}
	if !strings.Contains(err.Error(), "assinador.jar não encontrado") {
		t.Errorf("mensagem inesperada: %v", err)
	}
}

func TestFindJava_Ausente(t *testing.T) {
	isolarAmbiente(t)
	_, err := FindJava()
	if err == nil {
		t.Fatal("esperava erro quando o java não está disponível")
	}
	if !strings.Contains(err.Error(), "java não encontrado") {
		t.Errorf("mensagem inesperada: %v", err)
	}
}
