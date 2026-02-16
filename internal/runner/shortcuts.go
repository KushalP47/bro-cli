package runner

// Shortcut represents a command shortcut.
type Shortcut struct {
	Name        string `yaml:"name"`
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
}

// LoadShortcuts loads shortcuts from global and project-level YAML files.
func LoadShortcuts(globalPath, projectPath string) ([]Shortcut, error) {
	// TODO: load and merge shortcuts from both paths
	return nil, nil
}
