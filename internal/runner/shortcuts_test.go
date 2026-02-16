package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadShortcuts(t *testing.T) {
	tmpDir := t.TempDir()

	globalPath := filepath.Join(tmpDir, "commands.yaml")
	projectPath := filepath.Join(tmpDir, "project.yaml")

	globalContent := `shortcuts:
  - name: test
    command: go test ./...
    description: Run all tests
  - name: build
    command: go build -o app .
    description: Build the app
`
	projectContent := `shortcuts:
  - name: test
    command: npm test
    description: Run project tests
  - name: deploy
    command: ./deploy.sh
    description: Deploy to staging
`

	os.WriteFile(globalPath, []byte(globalContent), 0644)
	os.WriteFile(projectPath, []byte(projectContent), 0644)

	t.Run("merge with project override", func(t *testing.T) {
		shortcuts, err := LoadShortcuts(globalPath, projectPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(shortcuts) != 3 {
			t.Fatalf("expected 3 shortcuts, got %d", len(shortcuts))
		}

		// "test" should be overridden by project
		found := false
		for _, s := range shortcuts {
			if s.Name == "test" {
				found = true
				if s.Command != "npm test" {
					t.Errorf("expected project override, got command %q", s.Command)
				}
			}
		}
		if !found {
			t.Error("test shortcut not found")
		}
	})

	t.Run("global only", func(t *testing.T) {
		shortcuts, err := LoadShortcuts(globalPath, filepath.Join(tmpDir, "nonexistent.yaml"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(shortcuts) != 2 {
			t.Fatalf("expected 2 shortcuts, got %d", len(shortcuts))
		}
	})

	t.Run("no files", func(t *testing.T) {
		shortcuts, err := LoadShortcuts(
			filepath.Join(tmpDir, "nope1.yaml"),
			filepath.Join(tmpDir, "nope2.yaml"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(shortcuts) != 0 {
			t.Fatalf("expected 0 shortcuts, got %d", len(shortcuts))
		}
	})
}

func TestApplyPreset(t *testing.T) {
	tests := []struct {
		name    string
		command string
		preset  map[string]string
		want    string
	}{
		{
			name:    "single placeholder",
			command: "docker build -t {{image}}",
			preset:  map[string]string{"image": "myapp:latest"},
			want:    "docker build -t myapp:latest",
		},
		{
			name:    "multiple placeholders",
			command: "kubectl apply -n {{namespace}} -f {{file}}",
			preset:  map[string]string{"namespace": "production", "file": "deploy.yaml"},
			want:    "kubectl apply -n production -f deploy.yaml",
		},
		{
			name:    "no placeholders",
			command: "echo hello",
			preset:  map[string]string{"key": "value"},
			want:    "echo hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyPreset(tt.command, tt.preset)
			if result != tt.want {
				t.Errorf("expected %q, got %q", tt.want, result)
			}
		})
	}
}
