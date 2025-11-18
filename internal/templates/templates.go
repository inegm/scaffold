package templates

import (
	"embed"
	"text/template"
)

//go:embed files/* files/.github/workflows/*
var templateFS embed.FS

// ProjectType represents the type of project to scaffold
type ProjectType string

const (
	TypeCLI     ProjectType = "cli"
	TypeLibrary ProjectType = "library"
	TypeService ProjectType = "service"
)

// ProjectConfig holds the configuration for generating a project
type ProjectConfig struct {
	ProjectName string
	ModulePath  string
	Author      string
	License     string
	Type        ProjectType
}

// GetTemplate returns a parsed template by name
func GetTemplate(name string) (*template.Template, error) {
	content, err := templateFS.ReadFile("files/" + name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Parse(string(content))
}

// GetWorkflowFile returns the raw content of a workflow file
func GetWorkflowFile(name string) ([]byte, error) {
	return templateFS.ReadFile("files/" + name)
}

// DirectoryStructure returns the directories to create based on project type
func DirectoryStructure(projectType ProjectType) []string {
	switch projectType {
	case TypeLibrary:
		// Libraries focus on reusable code with documentation and examples
		return []string{
			"internal", // Private library code
			"pkg",      // Public library code (main focus)
			"cmd",      // Example CLI tools/utilities
			"scripts",  // Build/development scripts
			"test",     // Test apps and test data
			"docs",     // Documentation
			"examples", // Usage examples
		}
	case TypeCLI:
		// CLI applications with executables and optional public libraries
		return []string{
			"cmd",      // Main application(s)
			"internal", // Private application code
			"pkg",      // Public libraries (if exposing reusable components)
			"api",      // API definitions (if CLI interacts with APIs)
			"configs",  // Configuration templates
			"scripts",  // Build/install/analysis scripts
			"build",    // Packaging and CI files
			"test",     // External test apps and data
			"docs",     // Documentation
			"examples", // Usage examples
		}
	case TypeService:
		// Services with web/API components and full infrastructure
		return []string{
			"cmd",      // Service entry point(s)
			"internal", // Private service code
			"pkg",      // Shared libraries
			"api",      // API definitions (OpenAPI, protobuf, etc.)
			"web",      // Web application components
			"configs",  // Configuration files
			"scripts",  // Build/deployment scripts
			"build",    // Packaging and CI/CD
			"test",     // Test apps and data
			"docs",     // Documentation
			"examples", // Example API usage
		}
	default:
		// Default to CLI structure
		return []string{
			"cmd",
			"internal",
			"pkg",
			"api",
			"configs",
			"scripts",
			"build",
			"test",
			"docs",
			"examples",
		}
	}
}
