package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// StagedDiff returns the output of git diff --staged.
func StagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff --staged: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
