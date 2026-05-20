package cmd

import (
	"fmt"

	"github.com/jusouzn/assinatura/internal/jdk"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Provisiona o JRE 21 automaticamente em ~/.hubsaude/jdk",
	Long: `Detecta se o Java 21 está disponível no sistema ou em ~/.hubsaude/jdk.
Se não estiver disponível, baixa e instala o Eclipse Temurin JRE 21
automaticamente para a plataforma atual (Windows/Linux/macOS).`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(_ *cobra.Command, _ []string) error {
	if path, err := jdk.FindJava(); err == nil {
		fmt.Printf("Java já disponível: %s\n", path)
		return nil
	}

	fmt.Println("Java não encontrado. Iniciando provisionamento automático...")
	if err := jdk.Provision(); err != nil {
		return fmt.Errorf("falha no provisionamento: %w", err)
	}

	path, err := jdk.FindJava()
	if err != nil {
		return fmt.Errorf("JRE instalado mas não encontrado — verifique ~/.hubsaude/jdk/bin/java")
	}
	fmt.Printf("JRE instalado com sucesso: %s\n", path)
	return nil
}
