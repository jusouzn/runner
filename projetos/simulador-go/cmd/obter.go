package cmd

import (
	"fmt"

	"github.com/jusouzn/simulador/internal/jdk"
	"github.com/jusouzn/simulador/internal/release"
	"github.com/spf13/cobra"
)

var obterCmd = &cobra.Command{
	Use:   "obter",
	Short: "Baixa o simulador.jar e o JRE (se necessário) via GitHub Releases",
	Long: `Consulta o release.json do repositório da disciplina para verificar
a versão mais recente do simulador.jar. Baixa apenas se a versão local
estiver desatualizada ou ausente. Também provisiona o JRE 21 automaticamente.`,
	RunE: runObter,
}

func init() {
	rootCmd.AddCommand(obterCmd)
}

func runObter(_ *cobra.Command, _ []string) error {
	fmt.Println("Verificando versão mais recente do simulador...")

	info, err := release.Fetch()
	if err != nil {
		return fmt.Errorf("não foi possível verificar a versão remota: %w", err)
	}

	local := release.LocalVersion()
	if local != "" && local == info.Jar.Version {
		fmt.Printf("simulador.jar já está na versão mais recente (%s).\n", local)
	} else {
		if local != "" {
			fmt.Printf("Atualizando de v%s → v%s\n", local, info.Jar.Version)
		}
		if err := release.Download(info); err != nil {
			return err
		}
	}

	// Provisiona JRE se necessário
	if _, err := jdk.FindJava(); err != nil {
		fmt.Println("JRE não encontrado. Iniciando provisionamento...")
		jreURL, isZip := release.JREURLForPlatform(info)
		if jreURL == "" {
			return fmt.Errorf("release.json não contém URL do JRE para esta plataforma")
		}
		if err := jdk.ProvisionFromURLs(jreURL, isZip); err != nil {
			return fmt.Errorf("falha ao provisionar JRE: %w", err)
		}
		fmt.Println("JRE instalado com sucesso.")
	} else {
		fmt.Println("JRE já disponível.")
	}

	return nil
}
