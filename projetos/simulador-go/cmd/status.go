package cmd

import (
	"fmt"

	"github.com/jusouzn/simulador/internal/process"
	"github.com/jusouzn/simulador/internal/release"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibe o status atual do simulador",
	Long:  `Consulta o endpoint /api/info do simulador e exibe o resultado formatado.`,
	RunE:  runStatus,
}

var statusPorta int

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().IntVar(&statusPorta, "porta", process.DefaultPort, "Porta do simulador")
}

func runStatus(_ *cobra.Command, _ []string) error {
	local := release.LocalVersion()
	if local != "" {
		fmt.Printf("Versão instalada: %s\n", local)
	}

	if !process.IsRunning(statusPorta) {
		fmt.Printf("Simulador não está em execução (porta %d).\n", statusPorta)
		return nil
	}

	info, err := process.Status(statusPorta)
	if err != nil {
		return err
	}

	fmt.Printf("Simulador em execução na porta %d:\n%s\n", statusPorta, info)
	return nil
}
