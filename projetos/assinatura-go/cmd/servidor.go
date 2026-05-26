package cmd

import (
	"fmt"

	"github.com/jusouzn/assinatura/internal/server"
	"github.com/spf13/cobra"
)

var servidorCmd = &cobra.Command{
	Use:   "servidor",
	Short: "Controla o ciclo de vida do servidor assinador.jar",
}

var (
	pararPorta int
	pararApos  int
)

var servidorPararCmd = &cobra.Command{
	Use:   "parar",
	Short: "Encerra o servidor assinador.jar",
	Example: `  assinatura servidor parar
  assinatura servidor parar --porta 9090
  assinatura servidor parar --apos 10`,
	RunE: runServidorParar,
}

func init() {
	rootCmd.AddCommand(servidorCmd)
	servidorCmd.AddCommand(servidorPararCmd)

	servidorPararCmd.Flags().IntVar(&pararPorta, "porta", server.DefaultPort, "Porta do servidor assinador")
	servidorPararCmd.Flags().IntVar(&pararApos, "apos", 0, "Encerrar após N minutos sem interação (0 = imediato)")
}

func runServidorParar(_ *cobra.Command, _ []string) error {
	if err := server.StopOrSchedule(pararPorta, pararApos); err != nil {
		return err
	}
	if pararApos > 0 {
		fmt.Printf("Encerramento programado: servidor será parado após %d minuto(s) sem interação.\n", pararApos)
	} else {
		fmt.Println("Servidor assinador.jar encerrado.")
	}
	return nil
}
