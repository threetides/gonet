/*
Copyright © 2026 threetides.dev hello@threetides.dev
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gonet",
	Short: "Scaffold a Go HTTP service in seconds",
	Long: `gonet is a project scaffolder for Go HTTP services.

It generates a ready-to-run project: a standard folder layout, an
initialized Go module, an optional git repository, and a small httpx
package with a shared HTTP client and JSON helpers — so you can skip the
boilerplate and start writing handlers.

Get started with:

  gonet init my-service`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
