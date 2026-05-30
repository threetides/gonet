/*
Copyright © 2026 NAME HERE hello@threetides.dev
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func runCommand(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
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

func initProject(dir string) {
	// Create root dir
	err := os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created directory %s\n", dir)

	// go mod init
	var moduleName string
	scanner := bufio.NewScanner(os.Stdin)

	if dir == "." {
		fmt.Print("Module name: ")
	} else {
		fmt.Printf("Module name (press Enter to name it '%s'): ", dir)
	}
	if scanner.Scan() {
		if scanner.Text() == "" {
			moduleName = dir
		} else {
			moduleName = scanner.Text()
		}
	}

	runCommand(dir, "go", "mod", "init", moduleName)

	// Init git repo
	prompt := promptui.Select{
		Label: "Initialize git repository?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Fatal(err)
		return
	}

	if result == "Yes" {
		runCommand(dir, "git", "init")
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
				initProject("")
			} else {
				initProject(args[0])
			}
		} else {
			fmt.Print("Project name: ")
			_, err := fmt.Scanln(&dir)
			if err != nil {
				log.Fatal(err)
			}
			initProject(dir)
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
