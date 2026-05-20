package cmd

import (
	"fmt"

	"github.com/jusouzn/assinatura/internal/invoker"
	"github.com/jusouzn/assinatura/internal/server"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria uma assinatura digital (simulada) via assinador.jar",
	Example: `  assinatura sign --conteudo "texto a assinar" --pin 1234
  assinatura sign --conteudo "texto" --pin 1234 --algoritmo SHA512withRSA
  assinatura sign --local --conteudo "texto" --pin 1234`,
	RunE: runSign,
}

var (
	signConteudo  string
	signPin       string
	signAlgoritmo string
	signLocal     bool
)

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVar(&signConteudo, "conteudo", "", "Conteúdo a ser assinado (obrigatório)")
	signCmd.Flags().StringVar(&signPin, "pin", "", "PIN do dispositivo criptográfico (obrigatório)")
	signCmd.Flags().StringVar(&signAlgoritmo, "algoritmo", "SHA256withRSA", "Algoritmo de assinatura")
	signCmd.Flags().BoolVar(&signLocal, "local", false, "Invoca o assinador.jar diretamente, sem usar o servidor HTTP (cold start)")

	_ = signCmd.MarkFlagRequired("conteudo")
	_ = signCmd.MarkFlagRequired("pin")
}

func runSign(_ *cobra.Command, _ []string) error {
	p := invoker.Params{
		Operacao:  "assinar",
		Conteudo:  signConteudo,
		Pin:       signPin,
		Algoritmo: signAlgoritmo,
	}

	var (
		result string
		err    error
	)

	if signLocal {
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
