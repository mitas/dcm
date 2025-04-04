package formatter

import (
	"fmt"
	"strings"

	"github.com/mitas/dcm/internal/model"
)

// ANSI color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Formatter handles output formatting with colors and emojis
type Formatter struct{}

// NewFormatter creates a new formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatProjectList formats the list of projects for display
func (f *Formatter) FormatProjectList(projects []model.Project) string {
	if len(projects) == 0 {
		return fmt.Sprintf("%s❌ No Docker Compose projects found%s", ColorRed, ColorReset)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s📋 Found %s%d%s Docker Compose projects:%s\n",
		ColorBold, ColorGreen, len(projects), ColorReset, ColorReset))

	for i, project := range projects {
		sb.WriteString(fmt.Sprintf("%s📁 %d.%s %s%s%s (%s/%s)\n",
			ColorBlue, i+1, ColorReset,
			ColorBold, project.Name, ColorReset,
			project.Path, project.File))
	}

	return sb.String()
}

// FormatProjectStatus formats the status of a project
func (f *Formatter) FormatProjectStatus(projectName, projectPath string, running bool, services map[string]string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n%s=== Status of %s%s%s (%s) ===%s\n",
		ColorBold, ColorBlue, projectName, ColorReset, projectPath, ColorReset))

	if len(services) == 0 {
		sb.WriteString(fmt.Sprintf("%s🛑 Project is not running (no containers)%s\n", ColorYellow, ColorReset))
		return sb.String()
	}

	for service, status := range services {
		isRunning := strings.Contains(strings.ToLower(status), "up") || strings.Contains(strings.ToLower(status), "running")
		if isRunning {
			sb.WriteString(fmt.Sprintf("%s🟢 %s: %srunning%s (%s)\n", ColorGreen, service, ColorGreen, ColorReset, status))
		} else {
			sb.WriteString(fmt.Sprintf("%s🔴 %s: %sstopped%s (%s)\n", ColorRed, service, ColorRed, ColorReset, status))
		}
	}

	return sb.String()
}

// FormatActionResult formats the result of an action
func (f *Formatter) FormatActionResult(result model.Result) string {
	if result.Success {
		return fmt.Sprintf("%s✅ %s%s", ColorGreen, result.Message, ColorReset)
	}
	return fmt.Sprintf("%s❌ %s: %v%s", ColorRed, result.Project.Name, result.Error, ColorReset)
}

// FormatProjectNotFound formats a message when a project is not found
func (f *Formatter) FormatProjectNotFound(projectName string) string {
	return fmt.Sprintf("%s❓ Project %s%s%s not found%s",
		ColorYellow, ColorBold, projectName, ColorReset, ColorReset)
}

// FormatActionStart formats the start of an action
func (f *Formatter) FormatActionStart(actionName string, projectName string) string {
	return fmt.Sprintf("%s🔄 %s Docker Compose project: %s%s%s",
		ColorBold, actionName, ColorBlue, projectName, ColorReset)
}

// FormatNoProjectsFound formats a message when no projects are found
func (f *Formatter) FormatNoProjectsFound() string {
	return fmt.Sprintf("%s❌ No Docker Compose projects found%s", ColorRed, ColorReset)
}

// FormatManagedProjectsList formats the list of managed projects
func (f *Formatter) FormatManagedProjectsList(projects []model.ManagedProject) string {
	if len(projects) == 0 {
		return fmt.Sprintf("%s❌ No managed projects found%s", ColorYellow, ColorReset)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s📋 Managed Projects:%s\n", ColorBold, ColorReset))

	for i, p := range projects {
		sb.WriteString(fmt.Sprintf("%s📌 %d.%s %s%s%s (alias) -> %s%s%s (%s)\n",
			ColorBlue, i+1, ColorReset,
			ColorBold, p.Alias, ColorReset,
			ColorGreen, p.Project.Name, ColorReset,
			p.Project.Path))
	}

	return sb.String()
}
