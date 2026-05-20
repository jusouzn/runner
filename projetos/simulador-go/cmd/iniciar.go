package cmd

import (
	"fmt"

	"github.com/jusouzn/simulador/internal/jdk"
	"github.com/jusouzn/simulador/internal/process"
	"github.com/jusouzn/simulador/internal/release"
	"github.com/spf13/cobra"
)

var iniciarCmd = &cobra.Command{
	Use:   "iniciar",
	Short: "Inicia o simulador.jar",
	Long: `Verifica se a porta padrão (8443) está disponível e inicia o simulador.jar.
Execute "simulador obter" antes se o simulador.jar ainda não estiver instalado.`,
	RunE: runIniciar,
}

var iniciarPorta int

func init() {
	rootCmd.AddCommand(iniciarCmd)
	iniciarCmd.Flags().IntVar(&iniciarPorta, "porta", process.DefaultPort, "Porta do simulador")
}

func runIniciar(_ *cobra.Command, _ []string) error {
	if process.IsRunning(iniciarPorta) {
		fmt.Printf("Simulador já está em execução na porta %d.\n", iniciarPorta)
		return nil
	}

	if !process.IsPortFree(iniciarPorta) {
		return fmt.Errorf("porta %d está em uso por outro processo — libere-a antes de iniciar o simulador", iniciarPorta)
	}

	jarPath, err := release.JarPath()
	if err != nil {
		return err
	}

	javaPath, err := jdk.FindJava()
	if err != nil {
		return fmt.Errorf("%w\n  Dica: execute `simulador obter` para provisionar o JRE", err)
	}

	fmt.Printf("Iniciando simulador na porta %d...\n", iniciarPorta)
	if err := process.Start(javaPath, jarPath, iniciarPorta); err != nil {
		return err
	}

	fmt.Printf("Simulador iniciado com sucesso na porta %d.\n", iniciarPorta)
	return nil
}
