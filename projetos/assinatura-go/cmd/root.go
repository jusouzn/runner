package cmd

import (
	"os"

	"github.com/jusouzn/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

// Verbose indica se o modo verboso de logging está ativo.
// Quando ativo, os invokers (HTTP/local) imprimem em stderr detalhes
// sobre a URL do servidor ou o comando Java executado.
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para operações de assinatura digital via assinador.jar",
	Long: `Sistema Runner — CLI para criar e validar assinaturas digitais
sem necessidade de configurar o ambiente Java manualmente.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Propaga a flag --verbose para o pacote invoker.
		invoker.Verbose = Verbose
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// Não exibe o help completo quando ocorre erro de execução (ex: flag inválida)
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false,
		"exibe logs detalhados de execução (URL do servidor ou comando Java)")
}
