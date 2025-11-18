package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/daedalus/scaffold/internal/templates"
)

// Generator handles project generation
type Generator struct {
	config  *templates.ProjectConfig
	baseDir string
	dryRun  bool
}

// New creates a new Generator instance
func New(config *templates.ProjectConfig, baseDir string, dryRun bool) *Generator {
	return &Generator{
		config:  config,
		baseDir: baseDir,
		dryRun:  dryRun,
	}
}

// Generate creates the project structure
func (g *Generator) Generate() error {
	// Check if directory already exists
	projectPath := filepath.Join(g.baseDir, g.config.ProjectName)
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists", projectPath)
	}

	if g.dryRun {
		fmt.Println("Dry run mode - no files will be created")
		return g.preview()
	}

	// Create base project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create directory structure
	if err := g.createDirectories(projectPath); err != nil {
		return err
	}

	// Generate files from templates
	if err := g.generateFiles(projectPath); err != nil {
		return err
	}

	fmt.Printf("✓ Project %s created successfully at %s\n", g.config.ProjectName, projectPath)
	return nil
}

// createDirectories creates the project directory structure
func (g *Generator) createDirectories(projectPath string) error {
	dirs := templates.DirectoryStructure(g.config.Type)

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create cmd subdirectory for the main application
	cmdPath := filepath.Join(projectPath, "cmd", g.config.ProjectName)
	if err := os.MkdirAll(cmdPath, 0755); err != nil {
		return fmt.Errorf("failed to create cmd directory: %w", err)
	}

	// Create .github/workflows directory for GitHub Actions
	workflowsPath := filepath.Join(projectPath, ".github", "workflows")
	if err := os.MkdirAll(workflowsPath, 0755); err != nil {
		return fmt.Errorf("failed to create .github/workflows directory: %w", err)
	}

	return nil
}

// generateFiles creates files from templates
func (g *Generator) generateFiles(projectPath string) error {
	// Select Makefile template based on project type
	makefileTemplate := g.getMakefileTemplate()

	files := map[string]string{
		"README.md":  "README.md.tmpl",
		"Makefile":   makefileTemplate,
		".gitignore": "gitignore.tmpl",
		"go.mod":     "go.mod.tmpl",
	}

	// Add Dockerfile for service projects
	if g.config.Type == templates.TypeService {
		files["Dockerfile"] = "Dockerfile-service.tmpl"
	}

	for fileName, templateName := range files {
		if err := g.generateFile(projectPath, fileName, templateName); err != nil {
			return err
		}
	}

	// Generate main.go in cmd directory
	mainPath := filepath.Join("cmd", g.config.ProjectName, "main.go")
	if err := g.generateFile(projectPath, mainPath, "main.go.tmpl"); err != nil {
		return err
	}

	// Generate GitHub Actions workflow files
	if err := g.generateWorkflowFiles(projectPath); err != nil {
		return err
	}

	// Create placeholder .gitkeep files in empty directories
	if err := g.createGitkeepFiles(projectPath); err != nil {
		return err
	}

	return nil
}

// getMakefileTemplate returns the appropriate Makefile template based on project type
func (g *Generator) getMakefileTemplate() string {
	switch g.config.Type {
	case templates.TypeLibrary:
		return "Makefile-library.tmpl"
	case templates.TypeService:
		return "Makefile-service.tmpl"
	case templates.TypeCLI:
		return "Makefile-cli.tmpl"
	default:
		return "Makefile-cli.tmpl"
	}
}

// generateWorkflowFiles generates GitHub Actions workflow files
func (g *Generator) generateWorkflowFiles(projectPath string) error {
	// Common test workflow (includes lint as first job)
	commonWorkflows := map[string]string{
		".github/workflows/test.yml": ".github/workflows/test.yml",
	}

	// Generate common workflows
	for filePath, templatePath := range commonWorkflows {
		if err := g.copyWorkflowFile(projectPath, filePath, templatePath); err != nil {
			return err
		}
	}

	// Generate project-type-specific release workflow
	var releaseWorkflow string
	switch g.config.Type {
	case templates.TypeLibrary:
		releaseWorkflow = ".github/workflows/release-library.yml"
	case templates.TypeService:
		releaseWorkflow = ".github/workflows/release-service.yml"
	case templates.TypeCLI:
		releaseWorkflow = ".github/workflows/release-cli.yml"
	default:
		releaseWorkflow = ".github/workflows/release-cli.yml"
	}

	return g.generateFile(projectPath, ".github/workflows/release.yml", releaseWorkflow)
}

// copyWorkflowFile copies a workflow file without templating (for static files)
func (g *Generator) copyWorkflowFile(projectPath, destPath, srcPath string) error {
	// Read the workflow file from embedded templates
	content, err := templates.GetWorkflowFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read workflow file %s: %w", srcPath, err)
	}

	// Write to destination
	filePath := filepath.Join(projectPath, destPath)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write workflow file %s: %w", destPath, err)
	}

	return nil
}

// generateFile generates a single file from a template
func (g *Generator) generateFile(projectPath, fileName, templateName string) error {
	tmpl, err := templates.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	filePath := filepath.Join(projectPath, fileName)
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fileName, err)
	}

	return nil
}

// createGitkeepFiles creates .gitkeep files in empty directories
func (g *Generator) createGitkeepFiles(projectPath string) error {
	// Get the directory structure for this project type
	dirs := templates.DirectoryStructure(g.config.Type)

	// Directories that will have generated content (skip .gitkeep for these)
	skipDirs := map[string]bool{
		"cmd": true, // Will contain main.go
	}

	for _, dir := range dirs {
		if skipDirs[dir] {
			continue
		}

		gitkeepPath := filepath.Join(projectPath, dir, ".gitkeep")
		if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
			// Skip if directory doesn't exist or can't be written
			continue
		}
	}

	return nil
}

// preview shows what would be created without actually creating anything
func (g *Generator) preview() error {
	projectPath := filepath.Join(g.baseDir, g.config.ProjectName)

	fmt.Printf("\nProject: %s\n", g.config.ProjectName)
	fmt.Printf("Location: %s\n", projectPath)
	fmt.Printf("Module Path: %s\n", g.config.ModulePath)
	fmt.Printf("Type: %s\n\n", g.config.Type)

	fmt.Println("Directories to be created:")
	dirs := templates.DirectoryStructure(g.config.Type)
	for _, dir := range dirs {
		fmt.Printf("  %s/\n", dir)
	}
	fmt.Printf("  cmd/%s/\n", g.config.ProjectName)

	fmt.Println("\nFiles to be created:")
	files := []string{
		"README.md",
		"Makefile",
		".gitignore",
		"go.mod",
		fmt.Sprintf("cmd/%s/main.go", g.config.ProjectName),
	}

	// Add Dockerfile for service projects
	if g.config.Type == templates.TypeService {
		files = append(files, "Dockerfile")
	}

	for _, file := range files {
		fmt.Printf("  %s\n", file)
	}

	return nil
}

// ValidateProjectName checks if the project name is valid
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if strings.ContainsAny(name, " /\\:*?\"<>|") {
		return fmt.Errorf("project name contains invalid characters")
	}

	return nil
}
