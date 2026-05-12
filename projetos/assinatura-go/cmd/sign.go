package cmd

import (
	"fmt"

	"github.com/jusouzn/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria uma assinatura digital (simulada) via assinador.jar",
	Example: `  assinatura sign --conteudo "texto a assinar" --pin 1234
  assinatura sign --conteudo "texto" --pin 1234 --algoritmo SHA512withRSA`,
	RunE: runSign,
}

var (
	signConteudo  string
	signPin       string
	signAlgoritmo string
)

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVar(&signConteudo, "conteudo", "", "Conteúdo a ser assinado (obrigatório)")
	signCmd.Flags().StringVar(&signPin, "pin", "", "PIN do dispositivo criptográfico (obrigatório)")
	signCmd.Flags().StringVar(&signAlgoritmo, "algoritmo", "SHA256withRSA", "Algoritmo de assinatura")

	_ = signCmd.MarkFlagRequired("conteudo")
	_ = signCmd.MarkFlagRequired("pin")
}

func runSign(_ *cobra.Command, _ []string) error {
	result, err := invoker.Local(invoker.Params{
		Operacao:  "assinar",
		Conteudo:  signConteudo,
		Pin:       signPin,
		Algoritmo: signAlgoritmo,
	})
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
