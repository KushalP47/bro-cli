package config

import (
	"path/filepath"
)

// Profile defines commit behavior and is matched to directories via glob patterns.
type Profile struct {
	Name           string `yaml:"name" mapstructure:"name"`
	DirPattern     string `yaml:"dir_pattern" mapstructure:"dir_pattern"`
	CommitStyle    string `yaml:"commit_style" mapstructure:"commit_style"` // "template" or "conventional"
	Template       string `yaml:"template" mapstructure:"template"`
	BranchPattern  string `yaml:"branch_pattern" mapstructure:"branch_pattern"`
}

// MatchProfile returns the first profile whose DirPattern matches the given directory.
func MatchProfile(profiles []Profile, dir string) *Profile {
	for i := range profiles {
		matched, err := filepath.Match(profiles[i].DirPattern, dir)
		if err == nil && matched {
			return &profiles[i]
		}
	}
	return nil
}
