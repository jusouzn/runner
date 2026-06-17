package server

import (
	"strings"
	"testing"
)

func TestStop_SemInstancia(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	err := Stop(DefaultPort)
	if err == nil {
		t.Fatal("esperava erro ao parar sem instância em execução")
	}
	if !strings.Contains(err.Error(), "não está em execução") {
		t.Errorf("mensagem inesperada: %v", err)
	}
}

func TestStopOrSchedule_AgendarSemServidor(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	// Sem servidor vivo, o agendamento via HTTP não é possível e deve falhar claro.
	err := StopOrSchedule(DefaultPort, 5)
	if err == nil {
		t.Fatal("esperava erro ao agendar encerramento sem servidor acessível")
	}
	if !strings.Contains(err.Error(), "não está acessível via HTTP") {
		t.Errorf("mensagem inesperada: %v", err)
	}
}
