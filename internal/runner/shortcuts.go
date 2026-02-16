package runner

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Shortcut represents a command shortcut.
type Shortcut struct {
	Name        string                       `yaml:"name"`
	Command     string                       `yaml:"command"`
	Description string                       `yaml:"description"`
	Presets     map[string]map[string]string  `yaml:"presets"`
}

type shortcutsFile struct {
	Shortcuts []Shortcut `yaml:"shortcuts"`
}

// LoadShortcuts loads shortcuts from global and project-level YAML files.
// Project shortcuts override global shortcuts with the same name.
func LoadShortcuts(globalPath, projectPath string) ([]Shortcut, error) {
	globalShortcuts, err := loadFromFile(globalPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load global shortcuts: %w", err)
	}

	projectShortcuts, err := loadFromFile(projectPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load project shortcuts: %w", err)
	}

	// Merge: project overrides global by name
	merged := make(map[string]Shortcut)
	for _, s := range globalShortcuts {
		merged[s.Name] = s
	}
	for _, s := range projectShortcuts {
		merged[s.Name] = s
	}

	result := make([]Shortcut, 0, len(merged))
	// Preserve order: global first, then project-only
	seen := make(map[string]bool)
	for _, s := range globalShortcuts {
		result = append(result, merged[s.Name])
		seen[s.Name] = true
	}
	for _, s := range projectShortcuts {
		if !seen[s.Name] {
			result = append(result, s)
		}
	}

	return result, nil
}

func loadFromFile(path string) ([]Shortcut, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var file shortcutsFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse shortcuts YAML: %w", err)
	}
	return file.Shortcuts, nil
}

// ApplyPreset replaces {{placeholder}} tokens in a command with preset values.
func ApplyPreset(command string, preset map[string]string) string {
	for key, value := range preset {
		command = strings.ReplaceAll(command, "{{"+key+"}}", value)
	}
	return command
}
