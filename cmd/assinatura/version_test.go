package main

import (
	"os/exec"
	"strings"
	"testing"
)

// Este teste garante que o comando 'version' funciona e imprime 'dev'
func TestVersionCommand(t *testing.T) {
	// Preparamos o comando como se fôssemos o utilizador no terminal
	cmd := exec.Command("go", "run", ".", "version")
	
	// Executamos e capturamos a saída (o que aparece no ecrã)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	// Convertamos a saída para texto normal
	saidaTexto := string(output)

	// Verificamos se a saída contém a palavra "dev" (que é a versão padrão)
	if !strings.Contains(saidaTexto, "dev") {
		t.Errorf("Esperava encontrar a palavra 'dev' na versão, mas obtive: %s", saidaTexto)
	}
}