// main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

const VERSION = "0.1.0"

func main() {
	// Verifica se o usuário passou algum comando (ex: "assinar" ou "validar")
	if len(os.Args) < 2 {
		fmt.Println("Erro: Comando não informado. Use 'assinar' ou 'validar'.")
		os.Exit(1)
	}

	comando := os.Args[1]

	// Suporte para exibir a versão
	if comando == "--version" || comando == "-v" {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// Estrutura básica de comandos
	switch comando {
	case "assinar":
		chamarAssinadorLocal()
	default:
		fmt.Printf("Comando '%s' não reconhecido no modo walking skeleton.\n", comando)
		os.Exit(1)
	}
}

// Função responsável por invocar o assinador.jar localmente
func chamarAssinadorLocal() {
	fmt.Println("[CLI] Iniciando chamada para o assinador.jar em modo local...")

	// Prepara o comando: java -jar assinador.jar --modo=local
	// Nota: No futuro, o caminho do Java será resolvido dinamicamente.
	cmd := exec.Command("java", "-jar", "assinador.jar", "--modo=local", "--pin=1234")

	// Captura a saída padrão (stdout) e de erro (stderr) do processo Java
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("[CLI] Ocorreu um erro ao executar o JAR: %v\n", err)
		// Imprime a saída do erro em caso de falha
		fmt.Println(string(out))
		os.Exit(1)
	}

	// Imprime o resultado formatado (que deve ser o JSON vindo do Java)
	fmt.Println("[CLI] Resposta do Assinador:")
	fmt.Println(string(out))
}
