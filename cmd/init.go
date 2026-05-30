/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func makeDir(name string) {
	// Create root dir
	err := os.Mkdir(name, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created project %s\n", name)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string

		if len(args) == 1 {
			if args[0] == "." {
				fmt.Println("Making project in root")
			} else {
				makeDir(args[0])
			}
		} else {
			fmt.Print("Project name: ")
			_, err := fmt.Scanln(&projectName)
			if err != nil {
				log.Fatal(err)
			}
			makeDir(projectName)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
