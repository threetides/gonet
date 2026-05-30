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
	"path/filepath"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/threetides/gonet/internal/templates"
)

type Project struct {
	ModuleName string
}

func runCommand(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	if dir != "." {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalln("❌ Error running command:", err)
	}
}

func createFile(t string, filepath string, s any) {
	tmpl, err := template.ParseFS(templates.TemplateFS, t)
	if err != nil {
		log.Fatalln("Failed to parse template:", err)
	}

	// 3. Create or truncate the output file
	outputFile, err := os.Create(filepath)
	if err != nil {
		log.Fatalln("Failed to create output file:", err)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalln("Failed to close file:", err)
		}
	}()

	// 4. Render the template directly into the file
	err = tmpl.Execute(outputFile, s)
	if err != nil {
		log.Fatalln("Failed to execute template:", err)
	}

	fmt.Printf("✅ Created %s\n", filepath)
}

func initProject(projectName string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Ask for project name if not provided with args
	if len(projectName) == 0 {
		fmt.Print("Project name: ")
		_, err = fmt.Scanln(&projectName)
		if err != nil {
			log.Fatalln("❌ Error scanning response", err)
		}
	}

	// Create directories
	err = os.MkdirAll(filepath.Join(dir, projectName, "internal/httpx"), 0755)
	if err != nil {
		log.Fatalln("❌ Error creating directories:", err)
	}
	fmt.Println("✅ Folder structure initialized")

	// go mod init
	var moduleName string
	scanner := bufio.NewScanner(os.Stdin)

	if projectName == "." {
		fmt.Print("Module name: ")
	} else {
		fmt.Printf("Module name (press Enter to name it '%s'): ", projectName)
	}
	if scanner.Scan() {
		if scanner.Text() == "" {
			moduleName = projectName
		} else {
			moduleName = scanner.Text()
		}
	}

	runCommand(filepath.Join(dir, projectName), "go", "mod", "init", moduleName)
	runCommand(filepath.Join(dir, projectName), "go", "get", "github.com/joho/godotenv")

	// Init git repo
	prompt := promptui.Select{
		Label: "Initialize git repository?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Fatalln("❌ Error initializing git repo:", err)
		return
	}

	if result == "Yes" {
		runCommand(filepath.Join(dir, projectName), "git", "init")
	}

	// Create files and populate with templates
	createFile("main.go.tmpl", filepath.Join(dir, projectName, "/main.go"), Project{ModuleName: moduleName})
	createFile("helpers.go.tmpl", filepath.Join(dir, projectName, "/internal/httpx/helpers.go"), nil)
	createFile("types.go.tmpl", filepath.Join(dir, projectName, "/internal/httpx/types.go"), nil)
	createFile(".gitignore.tmpl", filepath.Join(dir, projectName, "/.gitignore"), nil)
	createFile("makefile.tmpl", filepath.Join(dir, projectName, "/makefile"), nil)
	_, err = os.Create(filepath.Join(dir, projectName, ".gitignore"))
	if err != nil {
		log.Fatalln("❌ Error creating .env:", err)
		return
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Short",
	Long:  `Long`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			initProject(args[0])
		} else {
			initProject("")
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
