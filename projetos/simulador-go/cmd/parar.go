package cmd

import (
	"fmt"

	"github.com/jusouzn/simulador/internal/process"
	"github.com/spf13/cobra"
)

var pararCmd = &cobra.Command{
	Use:   "parar",
	Short: "Para o simulador.jar em execução",
	Long: `Encerra o simulador via endpoint /shutdown. Se o endpoint não responder,
encerra o processo pelo PID salvo em ~/.hubsaude/simulador-server.json.`,
	RunE: runParar,
}

var pararPorta int

func init() {
	rootCmd.AddCommand(pararCmd)
	pararCmd.Flags().IntVar(&pararPorta, "porta", process.DefaultPort, "Porta do simulador")
}

func runParar(_ *cobra.Command, _ []string) error {
	if !process.IsRunning(pararPorta) {
		fmt.Printf("Simulador não está em execução na porta %d.\n", pararPorta)
		return nil
	}

	fmt.Printf("Encerrando simulador na porta %d...\n", pararPorta)
	if err := process.Stop(pararPorta); err != nil {
		return err
	}

	fmt.Println("Simulador encerrado.")
	return nil
}
