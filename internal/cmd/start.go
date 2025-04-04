package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mitas/dcm/internal/config"
	"github.com/mitas/dcm/internal/manager"
	"github.com/mitas/dcm/internal/model"
	"github.com/mitas/dcm/pkg/formatter"
)

// newStartCmd creates the start command
func newStartCmd(projectManager *manager.Manager, outputFormatter *formatter.Formatter) *cobra.Command {
	var all bool
	var projectName string

	cmd := &cobra.Command{
		Use:   "start [project]",
		Short: "Start docker-compose projects",
		Long:  `Start one or all docker-compose projects in the specified path or from managed projects.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no project name provided directly, check args
			if projectName == "" && len(args) > 0 {
				projectName = args[0]
			}

			// Check if it's a managed project first
			if projectName != "" && rootPath == "" {
				// Load managed config
				managedConfig, err := config.LoadManagedConfig(configPath)
				if err != nil {
					return fmt.Errorf("error loading managed projects: %w", err)
				}

				// Try to find it in managed projects
				managedProject, found := projectManager.FindManagedProject(managedConfig, projectName)
				if found {
					fmt.Println(outputFormatter.FormatActionStart("Starting managed", managedProject.Alias))
					result := projectManager.StartProject(managedProject.Project)
					fmt.Println(outputFormatter.FormatActionResult(result))
					return nil
				}

				// Project not found in managed projects
				if rootPath == "" {
					return fmt.Errorf("project '%s' not found in managed projects and no --path provided", projectName)
				}
			}

			// If we reach here, we need rootPath
			if rootPath == "" {
				return fmt.Errorf("path is required to find projects, use --path flag")
			}

			// Find all docker-compose projects
			projects, err := projectManager.FindProjects(rootPath)
			if err != nil {
				return fmt.Errorf("error finding projects: %w", err)
			}

			if len(projects) == 0 {
				fmt.Println(outputFormatter.FormatNoProjectsFound())
				return nil
			}

			// Set a timeout for the operation
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			if all {
				// Start all projects
				fmt.Printf("%s🔄 Starting %s%d%s Docker Compose projects...%s\n",
					formatter.ColorBold, formatter.ColorGreen, len(projects), formatter.ColorReset, formatter.ColorReset)

				results := projectManager.ManageAllProjects(ctx, projects, model.ActionStart)
				for _, result := range results {
					fmt.Println(outputFormatter.FormatActionResult(result))
				}
				return nil
			}

			// If no project name provided directly, check args
			if projectName == "" && len(args) > 0 {
				projectName = args[0]
			}

			// Validate project name
			if projectName == "" {
				return fmt.Errorf("project name is required when not using --all flag")
			}

			// Find and start the specific project
			project, found := projectManager.FindProject(projects, projectName)
			if !found {
				fmt.Println(outputFormatter.FormatProjectNotFound(projectName))
				return nil
			}

			fmt.Println(outputFormatter.FormatActionStart("Starting", project.Name))
			result := projectManager.StartProject(project)
			fmt.Println(outputFormatter.FormatActionResult(result))
			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Start all docker-compose projects")
	cmd.Flags().StringVarP(&projectName, "project", "n", "", "Name of the project to start")

	return cmd
}
