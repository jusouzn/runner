package cmd

import (
	"fmt"

	"github.com/jusouzn/assinatura/internal/invoker"
	"github.com/jusouzn/assinatura/internal/server"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida uma assinatura digital via assinador.jar",
	Example: `  assinatura validate --conteudo "texto original" --assinatura "MIIB..." --pin 1234
  assinatura validate --local --conteudo "texto" --assinatura "MIIB..." --pin 1234`,
	RunE: runValidate,
}

var (
	validateConteudo   string
	validatePin        string
	validateAssinatura string
	validateLocal      bool
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(&validateConteudo, "conteudo", "", "Conteúdo original (obrigatório)")
	validateCmd.Flags().StringVar(&validatePin, "pin", "", "PIN do dispositivo criptográfico (obrigatório)")
	validateCmd.Flags().StringVar(&validateAssinatura, "assinatura", "", "Assinatura a validar (obrigatório)")
	validateCmd.Flags().BoolVar(&validateLocal, "local", false, "Invoca o assinador.jar diretamente, sem usar o servidor HTTP (cold start)")

	_ = validateCmd.MarkFlagRequired("conteudo")
	_ = validateCmd.MarkFlagRequired("pin")
	_ = validateCmd.MarkFlagRequired("assinatura")
}

func runValidate(_ *cobra.Command, _ []string) error {
	p := invoker.Params{
		Operacao:   "validar",
		Conteudo:   validateConteudo,
		Pin:        validatePin,
		Assinatura: validateAssinatura,
	}

	var (
		result string
		err    error
	)

	if validateLocal {
		result, err = invoker.Local(p)
	} else {
		if err = server.EnsureRunning(server.DefaultPort); err != nil {
			return fmt.Errorf("não foi possível iniciar o servidor assinador: %w\n"+
				"  Dica: use --local para invocação direta (cold start)", err)
		}
		baseURL := fmt.Sprintf("http://localhost:%d", server.DefaultPort)
		result, err = invoker.HTTP(baseURL, p)
	}

	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
