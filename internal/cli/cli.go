package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mitas/dcm/internal/config"
	"github.com/mitas/dcm/internal/manager"
	"github.com/mitas/dcm/internal/model"
	"github.com/mitas/dcm/pkg/formatter"
)

// CLI handles command line interface operations
type CLI struct {
	manager   *manager.Manager
	formatter *formatter.Formatter
}

// NewCLI creates a new CLI
func NewCLI(manager *manager.Manager, formatter *formatter.Formatter) *CLI {
	return &CLI{
		manager:   manager,
		formatter: formatter,
	}
}

// Run executes the CLI with the given configuration
func (c *CLI) Run(cfg *config.Config) error {
	// Set a timeout for the entire operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Find all docker-compose projects
	projects, err := c.manager.FindProjects(cfg.RootPath)
	if err != nil {
		return fmt.Errorf("error finding projects: %w", err)
	}

	// Execute the requested action
	switch cfg.Action {
	case model.ActionList:
		fmt.Println(c.formatter.FormatProjectList(projects))

	case model.ActionStart:
		if cfg.ActionAll {
			// Start all projects
			results := c.manager.ManageAllProjects(ctx, projects, model.ActionStart)
			for _, result := range results {
				fmt.Println(c.formatter.FormatActionResult(result))
			}
		} else {
			// Start specific project
			project, found := c.manager.FindProject(projects, cfg.TargetProject)
			if !found {
				fmt.Println(c.formatter.FormatProjectNotFound(cfg.TargetProject))
				return nil
			}
			
			fmt.Println(c.formatter.FormatActionStart("Starting", project.Name))
			result := c.manager.StartProject(project)
			fmt.Println(c.formatter.FormatActionResult(result))
		}

	case model.ActionStop:
		if cfg.ActionAll {
			// Stop all projects
			results := c.manager.ManageAllProjects(ctx, projects, model.ActionStop)
			for _, result := range results {
				fmt.Println(c.formatter.FormatActionResult(result))
			}
		} else {
			// Stop specific project
			project, found := c.manager.FindProject(projects, cfg.TargetProject)
			if !found {
				fmt.Println(c.formatter.FormatProjectNotFound(cfg.TargetProject))
				return nil
			}
			
			fmt.Println(c.formatter.FormatActionStart("Stopping", project.Name))
			result := c.manager.StopProject(project)
			fmt.Println(c.formatter.FormatActionResult(result))
		}

	case model.ActionStatus:
		if cfg.ActionAll {
			// Check status of all projects
			fmt.Printf("🔍 Checking status of %d Docker Compose projects...\n", len(projects))
			for _, project := range projects {
				isRunning, services, err := c.manager.CheckProjectStatus(project)
				if err != nil {
					fmt.Printf("❌ Error checking status of %s: %v\n", project.Name, err)
					continue
				}
				fmt.Println(c.formatter.FormatProjectStatus(project.Name, project.Path, isRunning, services))
			}
		} else {
			// Check status of specific project
			project, found := c.manager.FindProject(projects, cfg.TargetProject)
			if !found {
				fmt.Println(c.formatter.FormatProjectNotFound(cfg.TargetProject))
				return nil
			}
			
			fmt.Println(c.formatter.FormatActionStart("Checking status of", project.Name))
			isRunning, services, err := c.manager.CheckProjectStatus(project)
			if err != nil {
				return fmt.Errorf("error checking status of %s: %w", project.Name, err)
			}
			fmt.Println(c.formatter.FormatProjectStatus(project.Name, project.Path, isRunning, services))
		}
	}

	return nil
}