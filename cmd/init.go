/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

func goModInit(name string) {
	fmt.Print("Module name: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "mod", "init", name)
	if name != "." {
		cmd.Dir = name
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Short",
	Long:  `Long`,
	Args:  cobra.MaximumNArgs(1),
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
			goModInit("")
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
