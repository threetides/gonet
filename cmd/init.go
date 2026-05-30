/*
Copyright © 2026 threetides.dev hello@threetides.dev
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

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/threetides/gonet/internal/templates"
)

type Project struct {
	ModuleName string
}

var red = color.New(color.FgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()

func runCommand(dir string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	if dir != "." {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalln(red("error:"), cyan("error running command:"), err)
	}
}

func createFile(t string, filepath string, s any) {
	tmpl, err := template.ParseFS(templates.TemplateFS, t)
	if err != nil {
		log.Fatalln(red("error:"), cyan("failed to parse template:"), err)
	}

	// 3. Create or truncate the output file
	outputFile, err := os.Create(filepath)
	if err != nil {
		log.Fatalln(red("error:"), cyan("failed to create output file:"), err)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalln(red("error:"), cyan("failed to close file:"), err)
		}
	}()

	// 4. Render the template directly into the file
	err = tmpl.Execute(outputFile, s)
	if err != nil {
		log.Fatalln(red("error:"), cyan("failed to execute template:"), err)
	}

	fmt.Printf("gonet: created %s\n", filepath)
}

func initProject(projectName string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(red("error:"), cyan("failed to get current path:"), err)
	}

	// Ask for project name if not provided with args
	if len(projectName) == 0 {
		fmt.Print(blue("enter project name:"))
		_, err = fmt.Scanln(&projectName)
		if err != nil {
			log.Fatalln(red("error:"), cyan("error scanning response:"), err)
		}
	}

	// Create directories
	err = os.MkdirAll(filepath.Join(dir, projectName, "internal/httpx"), 0755)
	if err != nil {
		log.Fatalln(red("error:"), cyan("error creating directories:"), err)
	}
	fmt.Println("gonet: folder structure initialized")

	// go mod init
	var moduleName string
	scanner := bufio.NewScanner(os.Stdin)

	if projectName == "." {
		fmt.Print("enter module name:")
	} else {
		fmt.Printf("enter module name (press enter to leave it as '%s'):", projectName)
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
		Label: "initialize a git repository?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Fatalln(red("error:"), cyan("error initializing git repo:"), err)
		return
	}

	if result == "yes" {
		runCommand(filepath.Join(dir, projectName), "git", "init")
	}

	// Create files and populate with templates
	createFile("main.go.tmpl", filepath.Join(dir, projectName, "/main.go"), Project{ModuleName: moduleName})
	createFile("helpers.go.tmpl", filepath.Join(dir, projectName, "/internal/httpx/helpers.go"), nil)
	createFile("types.go.tmpl", filepath.Join(dir, projectName, "/internal/httpx/types.go"), nil)
	createFile("client.go.tmpl", filepath.Join(dir, projectName, "/internal/httpx/client.go"), nil)
	createFile(".gitignore.tmpl", filepath.Join(dir, projectName, "/.gitignore"), nil)
	createFile("makefile.tmpl", filepath.Join(dir, projectName, "/makefile"), nil)
	_, err = os.Create(filepath.Join(dir, projectName, ".env"))
	if err != nil {
		log.Fatalln(red("error:"), cyan("error creating .env:"), err)
		return
	}

	fmt.Println(green("success:"), cyan("successfully initialized ", projectName, "!"))

	if projectName == "" {
		fmt.Println("gonet:", cyan("run 'go run main.go' or 'make dev' to get started"))
	} else {
		fmt.Println("gonet:", cyan("run cd"), cyan(projectName), cyan("&& 'go run main.go' or 'make dev' to get started"))
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Create a new Go HTTP project",
	Long: `Create a new Go HTTP project in a new directory.

Pass a project name to create it in a directory of that name, or use "."
to scaffold into the current directory. If no name is given, you will be
prompted for one.

The command creates the folder structure, initializes a Go module,
optionally sets up a git repository, and generates the starter files
(main.go, the internal/httpx package, makefile, .gitignore and .env).

Examples:
  gonet init my-service   # create ./my-service
  gonet init .            # scaffold into the current directory
  gonet init              # prompt for a name`,
	Args: cobra.MaximumNArgs(1),
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
}
