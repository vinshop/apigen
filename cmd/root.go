package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vinshop/apigen/cmd/modelgen"
	"go.uber.org/zap"
	"os"
)

var rootCmd *cobra.Command

func init() {

	rootCmd = &cobra.Command{
		Use:   "version",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	rootCmd.AddCommand(modelgen.GenAPI)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.S().Error("Error when execute command, detail: ", err)
		os.Exit(0)
	}
}
