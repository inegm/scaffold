package main

import (
	"fmt"
	"os"

	"github.com/daedalus/scaffold/internal/config"
	"github.com/daedalus/scaffold/internal/generator"
	"github.com/daedalus/scaffold/internal/templates"
	"github.com/spf13/cobra"
)

var (
	modulePath  string
	author      string
	license     string
	projectType string
	dryRun      bool
	interactive bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "A project scaffolding tool for Go projects",
	Long: `Scaffold is a CLI tool for creating new Go projects following the
golang-standards/project-layout structure.

It helps you quickly set up a well-structured Go project with all the
standard directories and files needed to get started.`,
	Version: "1.0.0",
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Long: `Create a new Go project with the standard project layout.

This command generates a new project directory with:
  - Standard directory structure (/cmd, /internal, /pkg, etc.)
  - Basic configuration files (README.md, Makefile, .gitignore)
  - Go module initialization
  - Sample main.go file

If no project name is provided, interactive mode will be used.

Examples:
  scaffold new                  (interactive mode)
  scaffold new myapp
  scaffold new myapp --type cli --module-path github.com/user/myapp
  scaffold new myapp --dry-run`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNew,
}

func init() {
	newCmd.Flags().StringVarP(&modulePath, "module-path", "m", "", "Go module path (e.g., github.com/user/project)")
	newCmd.Flags().StringVarP(&author, "author", "a", "", "Author name")
	newCmd.Flags().StringVarP(&license, "license", "l", "MIT", "License type")
	newCmd.Flags().StringVarP(&projectType, "type", "t", "cli", "Project type (cli, library, service)")
	newCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview what would be created without creating anything")
	newCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Use interactive mode to configure the project")

	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	var projectName string

	// If no project name provided, enable interactive mode
	if len(args) == 0 {
		interactive = true
		projectName = "" // Will be prompted for
	} else {
		projectName = args[0]

		// Validate project name if provided
		if err := generator.ValidateProjectName(projectName); err != nil {
			return err
		}
	}

	// Create initial project config
	pType := templates.ProjectType(projectType)
	cfg := &templates.ProjectConfig{
		ProjectName: projectName,
		ModulePath:  modulePath,
		Author:      author,
		License:     license,
		Type:        pType,
	}

	// Use interactive mode if explicitly requested or if no project name was provided
	if interactive {
		if err := config.PromptForMissingConfig(cfg); err != nil {
			return err
		}

		// Validate the project name after interactive input
		if err := generator.ValidateProjectName(cfg.ProjectName); err != nil {
			return err
		}
	} else {
		// Validate project type when not in interactive mode
		if pType != templates.TypeCLI && pType != templates.TypeLibrary && pType != templates.TypeService {
			return fmt.Errorf("invalid project type: %s (must be cli, library, or service)", projectType)
		}

		// Set defaults for non-interactive mode
		if cfg.ModulePath == "" {
			cfg.ModulePath = projectName
		}
		if cfg.Author == "" {
			cfg.Author = "Your Name"
		}
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create generator and generate project
	gen := generator.New(cfg, cwd, dryRun)
	if err := gen.Generate(); err != nil {
		return err
	}

	if !dryRun {
		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  make deps")
		fmt.Println("  make build")
		fmt.Println("  make run")
	}

	return nil
}
