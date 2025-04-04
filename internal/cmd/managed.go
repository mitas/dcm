package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mitas/dcm/internal/config"
	"github.com/mitas/dcm/internal/manager"
	"github.com/mitas/dcm/pkg/formatter"
)

// newListManagedCmd creates a command to list managed projects
func newListManagedCmd(projectManager *manager.Manager, outputFormatter *formatter.Formatter) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-managed",
		Aliases: []string{"lsm"},
		Short:   "List all managed docker-compose projects",
		Long:    `List all docker-compose projects that have been saved to the config file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load managed config
			managedConfig, err := config.LoadManagedConfig(configPath)
			if err != nil {
				return fmt.Errorf("error loading managed projects: %w", err)
			}

			// Format and display managed projects
			if len(managedConfig.Projects) == 0 {
				fmt.Printf("%sNo managed projects found%s\n", formatter.ColorYellow, formatter.ColorReset)
				return nil
			}

			fmt.Printf("%s📋 Managed Projects:%s\n", formatter.ColorBold, formatter.ColorReset)
			for i, p := range managedConfig.Projects {
				fmt.Printf("%s📌 %d.%s %s%s%s (alias) -> %s%s%s (%s)\n",
					formatter.ColorBlue, i+1, formatter.ColorReset,
					formatter.ColorBold, p.Alias, formatter.ColorReset,
					formatter.ColorGreen, p.Project.Name, formatter.ColorReset,
					p.Project.Path)
			}

			return nil
		},
	}

	return cmd
}

// newAddManagedCmd creates a command to add a managed project
func newAddManagedCmd(projectManager *manager.Manager, outputFormatter *formatter.Formatter) *cobra.Command {
	var alias string

	cmd := &cobra.Command{
		Use:     "add-managed [project]",
		Aliases: []string{"add"},
		Short:   "Add a project to managed projects",
		Long:    `Add a docker-compose project to the managed projects list with an optional alias.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if rootPath == "" {
				return fmt.Errorf("path is required to find projects, use --path flag")
			}

			if len(args) < 1 {
				return fmt.Errorf("project name is required")
			}

			projectName := args[0]

			// Find projects in the specified path
			projects, err := projectManager.FindProjects(rootPath)
			if err != nil {
				return fmt.Errorf("error finding projects: %w", err)
			}

			// Find the target project
			project, found := projectManager.FindProject(projects, projectName)
			if !found {
				return fmt.Errorf("project '%s' not found in %s", projectName, rootPath)
			}

			// If alias is not provided, use the project name
			if alias == "" {
				alias = project.Name
			}

			// Load managed config
			managedConfig, err := config.LoadManagedConfig(configPath)
			if err != nil {
				return fmt.Errorf("error loading managed projects: %w", err)
			}

			// Add the project to managed projects
			if err := projectManager.AddManagedProject(managedConfig, project, alias); err != nil {
				return fmt.Errorf("error adding managed project: %w", err)
			}

			// Save managed config
			if err := config.SaveManagedConfig(managedConfig, configPath); err != nil {
				return fmt.Errorf("error saving managed projects: %w", err)
			}

			fmt.Printf("%s✅ Project '%s%s%s' added to managed projects with alias '%s%s%s'%s\n",
				formatter.ColorGreen,
				formatter.ColorBold, project.Name, formatter.ColorReset+formatter.ColorGreen,
				formatter.ColorBold, alias, formatter.ColorReset+formatter.ColorGreen,
				formatter.ColorReset)
			return nil
		},
	}

	cmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias for the managed project (defaults to project name)")

	return cmd
}

// newRemoveManagedCmd creates a command to remove a managed project
func newRemoveManagedCmd(projectManager *manager.Manager, outputFormatter *formatter.Formatter) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove-managed [alias]",
		Aliases: []string{"rm"},
		Short:   "Remove a project from managed projects",
		Long:    `Remove a docker-compose project from the managed projects list using its alias.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("project alias is required")
			}

			alias := args[0]

			// Load managed config
			managedConfig, err := config.LoadManagedConfig(configPath)
			if err != nil {
				return fmt.Errorf("error loading managed projects: %w", err)
			}

			// Remove the project from managed projects
			if err := projectManager.RemoveManagedProject(managedConfig, alias); err != nil {
				return fmt.Errorf("error removing managed project: %w", err)
			}

			// Save managed config
			if err := config.SaveManagedConfig(managedConfig, configPath); err != nil {
				return fmt.Errorf("error saving managed projects: %w", err)
			}

			fmt.Printf("%s✅ Project with alias '%s%s%s' removed from managed projects%s\n",
				formatter.ColorGreen,
				formatter.ColorBold, alias, formatter.ColorReset+formatter.ColorGreen,
				formatter.ColorReset)
			return nil
		},
	}

	return cmd
}