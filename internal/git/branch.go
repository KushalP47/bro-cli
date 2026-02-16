package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// CurrentBranch returns the current git branch name.
func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git current branch: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ParseBranch extracts named groups from a branch name using the given regex pattern.
func ParseBranch(branch, pattern string) (map[string]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("compile branch pattern: %w", err)
	}

	match := re.FindStringSubmatch(branch)
	if match == nil {
		return nil, fmt.Errorf("branch %q does not match pattern %q", branch, pattern)
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, nil
}
