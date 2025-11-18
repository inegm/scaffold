package config

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/daedalus/scaffold/internal/templates"
)

// PromptForMissingConfig prompts the user for any missing configuration
func PromptForMissingConfig(config *templates.ProjectConfig) error {
	var projectType string

	// Convert ProjectType to string for the form
	switch config.Type {
	case templates.TypeCLI:
		projectType = "cli"
	case templates.TypeLibrary:
		projectType = "library"
	case templates.TypeService:
		projectType = "service"
	default:
		projectType = "cli"
	}

	// Build form fields - add project name field if not already set
	var fields []huh.Field

	if config.ProjectName == "" {
		fields = append(fields,
			huh.NewInput().
				Title("Project Name").
				Description("The name of your project (used for directory and binary name)").
				Value(&config.ProjectName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("project name cannot be empty")
					}
					return nil
				}),
		)
	}

	fields = append(fields,
		huh.NewInput().
			Title("Module Path").
			Description("The Go module path (e.g., github.com/user/project)").
			Value(&config.ModulePath).
			Placeholder(func() string {
				if config.ProjectName != "" {
					return config.ProjectName
				}
				return "myproject"
			}()),

		huh.NewInput().
			Title("Author").
			Description("Your name or organization").
			Value(&config.Author).
			Placeholder("Your Name"),

		huh.NewSelect[string]().
			Title("License").
			Description("Choose a license for your project").
			Options(
				huh.NewOption("MIT", "MIT"),
				huh.NewOption("Apache 2.0", "Apache-2.0"),
				huh.NewOption("GPL v3", "GPL-3.0"),
				huh.NewOption("BSD 3-Clause", "BSD-3-Clause"),
				huh.NewOption("None", "None"),
			).
			Value(&config.License),

		huh.NewSelect[string]().
			Title("Project Type").
			Description("What type of project are you building?").
			Options(
				huh.NewOption("CLI Application", "cli"),
				huh.NewOption("Library", "library"),
				huh.NewOption("Service/API", "service"),
			).
			Value(&projectType),
	)

	form := huh.NewForm(
		huh.NewGroup(fields...),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form cancelled or failed: %w", err)
	}

	// Update config with form values
	config.Type = templates.ProjectType(projectType)

	// Set default module path if still empty
	if config.ModulePath == "" {
		config.ModulePath = config.ProjectName
	}

	// Set default author if still empty
	if config.Author == "" {
		config.Author = "Your Name"
	}

	return nil
}
