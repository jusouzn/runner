package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version é injetada via ldflags durante o build de release: -X github.com/jusouzn/assinatura/cmd.Version=X.Y.Z
var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("assinatura v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
