package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// A versão é declarada como variável (e não constante) para que o 
// GitHub Actions consiga injetar a versão real (ex: v0.1.0) durante a compilação.
var version = "dev"

func main() {
	// 1. Definimos o comando principal (a "raiz" do nosso CLI)
	var rootCmd = &cobra.Command{
		Use:   "assinatura",
		Short: "CLI para o Sistema Runner",
		Long:  `Ferramenta de linha de comandos para criar e validar assinaturas digitais simuladas.`,
	}

	// 2. Definimos o subcomando "version"
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual do CLI",
		Run: func(cmd *cobra.Command, args []string) {
			// runtime.GOOS e runtime.GOARCH mostram o sistema operativo e arquitetura
			fmt.Printf("Assinatura CLI v%s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}

	// 3. Adicionamos o subcomando ao comando principal
	rootCmd.AddCommand(versionCmd)

	// 4. Executamos o CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}