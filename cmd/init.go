/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
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

func goModInit(dir string) {
	var name string
	scanner := bufio.NewScanner(os.Stdin)

	if dir == "." {
		fmt.Print("Module name: ")
	} else {
		fmt.Printf("Module name (press Enter to name it '%s'): ", dir)
	}
	if scanner.Scan() {
		if scanner.Text() == "" {
			name = dir
		} else {
			name = scanner.Text()
		}
	}

	cmd := exec.Command("go", "mod", "init", name)
	if dir != "." {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
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
		var dir string

		if len(args) == 1 {
			if args[0] == "." {
				fmt.Println("Making project in root")
				goModInit("")
			} else {
				makeDir(args[0])
				goModInit(args[0])
			}
		} else {
			fmt.Print("Project name: ")
			_, err := fmt.Scanln(&dir)
			if err != nil {
				log.Fatal(err)
			}
			makeDir(dir)
			goModInit(dir)
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
