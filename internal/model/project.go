package model

// Project represents a docker-compose project
type Project struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	File string `yaml:"file"`
}

// ManagedProject represents a saved docker-compose project
type ManagedProject struct {
	// Alias is a user-friendly name for the project
	Alias string `yaml:"alias"`
	// Project contains the actual project data
	Project Project `yaml:"project"`
}

// ActionType defines what action to perform on docker-compose projects
type ActionType int

const (
	// ActionList lists all projects
	ActionList ActionType = iota
	// ActionStart starts a project
	ActionStart
	// ActionStop stops a project
	ActionStop
	// ActionStatus checks status of a project
	ActionStatus
)

// Result represents the result of a docker-compose operation
type Result struct {
	Project Project
	Success bool
	Message string
	Error   error
}