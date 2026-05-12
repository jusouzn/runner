package cmd

import (
	"fmt"

	"github.com/jusouzn/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida uma assinatura digital via assinador.jar",
	Example: `  assinatura validate --conteudo "texto original" --assinatura "MIIB..." --pin 1234`,
	RunE: runValidate,
}

var (
	validateConteudo   string
	validatePin        string
	validateAssinatura string
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(&validateConteudo, "conteudo", "", "Conteúdo original (obrigatório)")
	validateCmd.Flags().StringVar(&validatePin, "pin", "", "PIN do dispositivo criptográfico (obrigatório)")
	validateCmd.Flags().StringVar(&validateAssinatura, "assinatura", "", "Assinatura a validar (obrigatório)")

	_ = validateCmd.MarkFlagRequired("conteudo")
	_ = validateCmd.MarkFlagRequired("pin")
	_ = validateCmd.MarkFlagRequired("assinatura")
}

func runValidate(_ *cobra.Command, _ []string) error {
	result, err := invoker.Local(invoker.Params{
		Operacao:   "validar",
		Conteudo:   validateConteudo,
		Pin:        validatePin,
		Assinatura: validateAssinatura,
	})
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
