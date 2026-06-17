package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Variáveis injetadas via ldflags no build de release, ex.:
//
//	-X github.com/jusouzn/assinatura/cmd.Version=1.2.3
//	-X github.com/jusouzn/assinatura/cmd.Commit=abc1234
//	-X github.com/jusouzn/assinatura/cmd.Date=2026-06-17
var (
	Version = "0.1.0-dev"
	Commit  = ""
	Date    = ""
)

// versionString monta uma linha de versão rastreável: tag + SHA curto (+ data).
func versionString() string {
	s := "assinatura v" + Version
	if Commit != "" {
		s += " (" + Commit
		if Date != "" {
			s += ", " + Date
		}
		s += ")"
	}
	return s
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do CLI (tag + commit)",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(versionString())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
