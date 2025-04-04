package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/mitas/dcm/internal/config"
	"github.com/mitas/dcm/internal/model"
)

// CommandExecutor executes shell commands
type CommandExecutor interface {
	Execute(dir string, command string, args ...string) ([]byte, error)
}

// DefaultCommandExecutor is the default implementation of CommandExecutor
type DefaultCommandExecutor struct{}

// Execute runs a command and returns its output
func (e *DefaultCommandExecutor) Execute(dir string, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// Manager handles docker-compose operations
type Manager struct {
	executor CommandExecutor
}

// NewManager creates a new manager
func NewManager(executor CommandExecutor) *Manager {
	if executor == nil {
		executor = &DefaultCommandExecutor{}
	}
	return &Manager{
		executor: executor,
	}
}

// FindProjects searches for docker-compose projects in the given path
func (m *Manager) FindProjects(rootPath string) ([]model.Project, error) {
	var projects []model.Project

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip permission errors and continue with the walk
			return nil
		}

		// Skip directories that start with . (hidden directories) or # (temp files)
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(info.Name(), "#")) {
			return filepath.SkipDir
		}

		// Check if the file is a docker-compose file
		filename := info.Name()
		if !info.IsDir() && (filename == "docker-compose.yml" || filename == "docker-compose.yaml") {
			dirPath := filepath.Dir(path)
			projectName := filepath.Base(dirPath)
			projects = append(projects, model.Project{
				Name: projectName,
				Path: dirPath,
				File: filename,
			})
		}
		return nil
	})

	return projects, err
}

// StartProject starts a docker-compose project
func (m *Manager) StartProject(project model.Project) model.Result {
	output, err := m.executor.Execute(project.Path, "docker", "compose", "up", "-d")
	if err != nil {
		return model.Result{
			Project: project,
			Success: false,
			Error:   fmt.Errorf("error starting %s: %w: %s", project.Name, err, output),
		}
	}
	return model.Result{
		Project: project,
		Success: true,
		Message: fmt.Sprintf("Successfully started %s", project.Name),
	}
}

// StopProject stops a docker-compose project
func (m *Manager) StopProject(project model.Project) model.Result {
	output, err := m.executor.Execute(project.Path, "docker", "compose", "down")
	if err != nil {
		return model.Result{
			Project: project,
			Success: false,
			Error:   fmt.Errorf("error stopping %s: %w: %s", project.Name, err, output),
		}
	}
	return model.Result{
		Project: project,
		Success: true,
		Message: fmt.Sprintf("Successfully stopped %s", project.Name),
	}
}

// CheckProjectStatus checks the status of a docker-compose project
func (m *Manager) CheckProjectStatus(project model.Project) (bool, map[string]string, error) {
	// Check if any containers exist
	output, err := m.executor.Execute(project.Path, "docker", "compose", "ps", "-a", "--format", "json")
	if err != nil {
		return false, nil, fmt.Errorf("error checking status: %w", err)
	}

	// If no output or just whitespace, project is not running
	if len(strings.TrimSpace(string(output))) == 0 || string(output) == "[]" {
		return false, nil, nil
	}

	// Get services from docker-compose.yml
	servicesOutput, err := m.executor.Execute(project.Path, "docker", "compose", "config", "--services")
	if err != nil {
		return false, nil, fmt.Errorf("error getting services: %w", err)
	}

	services := strings.Split(strings.TrimSpace(string(servicesOutput)), "\n")
	if len(services) == 0 || (len(services) == 1 && services[0] == "") {
		return false, nil, nil
	}

	// Get status for each service
	serviceStatus := make(map[string]string)
	isRunning := false

	for _, service := range services {
		if service == "" {
			continue
		}

		statusOutput, err := m.executor.Execute(project.Path, "docker", "compose", "ps", service, "--format", "{{.Status}}")
		if err != nil || strings.TrimSpace(string(statusOutput)) == "" {
			serviceStatus[service] = "not running"
			continue
		}

		status := strings.TrimSpace(string(statusOutput))
		serviceStatus[service] = status

		// Check if service is running
		if regexp.MustCompile(`(?i)up|running`).MatchString(status) {
			isRunning = true
		}
	}

	return isRunning, serviceStatus, nil
}

// ManageAllProjects executes an action on all projects concurrently
func (m *Manager) ManageAllProjects(ctx context.Context, projects []model.Project, action model.ActionType) []model.Result {
	if len(projects) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	resultCh := make(chan model.Result, len(projects))
	results := make([]model.Result, 0, len(projects))

	// Process each project in a goroutine
	for _, project := range projects {
		wg.Add(1)
		go func(p model.Project) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				resultCh <- model.Result{
					Project: p,
					Success: false,
					Error:   ctx.Err(),
				}
				return
			default:
				var result model.Result
				switch action {
				case model.ActionStart:
					result = m.StartProject(p)
				case model.ActionStop:
					result = m.StopProject(p)
				}
				resultCh <- result
			}
		}(project)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	for result := range resultCh {
		results = append(results, result)
	}

	return results
}

// FindProject finds a project by name
func (m *Manager) FindProject(projects []model.Project, projectName string) (model.Project, bool) {
	for _, project := range projects {
		if strings.Contains(strings.ToLower(project.Name), strings.ToLower(projectName)) {
			return project, true
		}
	}
	return model.Project{}, false
}

// AddManagedProject adds a project to the managed projects list
func (m *Manager) AddManagedProject(managedConfig *config.ManagedConfig, project model.Project, alias string) error {
	// Check if alias already exists
	for _, p := range managedConfig.Projects {
		if p.Alias == alias {
			return fmt.Errorf("a project with alias '%s' already exists", alias)
		}
	}

	// Add the project to config
	managedConfig.Projects = append(managedConfig.Projects, model.ManagedProject{
		Alias:   alias,
		Project: project,
	})

	return nil
}

// RemoveManagedProject removes a project from the managed projects list
func (m *Manager) RemoveManagedProject(managedConfig *config.ManagedConfig, alias string) error {
	for i, p := range managedConfig.Projects {
		if p.Alias == alias {
			// Remove project at index i
			managedConfig.Projects = append(managedConfig.Projects[:i], managedConfig.Projects[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("no project found with alias '%s'", alias)
}

// FindManagedProject finds a managed project by alias
func (m *Manager) FindManagedProject(managedConfig *config.ManagedConfig, alias string) (model.ManagedProject, bool) {
	for _, p := range managedConfig.Projects {
		if strings.EqualFold(p.Alias, alias) || strings.Contains(strings.ToLower(p.Alias), strings.ToLower(alias)) {
			return p, true
		}
	}
	return model.ManagedProject{}, false
}