package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mitas/dcm/internal/manager"
	"github.com/mitas/dcm/pkg/formatter"
)

// newListCmd creates the list command
func newListCmd(projectManager *manager.Manager, outputFormatter *formatter.Formatter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all docker-compose projects",
		Long:  `Find and list all docker-compose projects in the specified path.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find all docker-compose projects
			projects, err := projectManager.FindProjects(rootPath)
			if err != nil {
				return fmt.Errorf("error finding projects: %w", err)
			}

			// Print the formatted list
			fmt.Println(outputFormatter.FormatProjectList(projects))
			return nil
		},
	}

	return cmd
}
