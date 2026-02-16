package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Log returns git log entries since the given date for a repo path.
func Log(repoPath, since string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--since="+since, "--oneline", "--no-merges")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log for %s: %w", repoPath, err)
	}
	return strings.TrimSpace(string(out)), nil
}
