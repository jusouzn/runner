package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "simulador",
	Short: "CLI para gerenciar o Simulador do HubSaúde",
	Long: `Sistema Runner — CLI para iniciar, parar e monitorar o simulador.jar
do HubSaúde, sem necessidade de conhecer os comandos Java subjacentes.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
}
