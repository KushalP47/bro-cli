package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/kushalpatel/bro-cli/internal/config"
	"github.com/kushalpatel/bro-cli/internal/runner"
	"github.com/kushalpatel/bro-cli/internal/ui"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [alias]",
	Short: "Run a command shortcut",
	Long:  "Fuzzy-searchable command shortcuts with arrow key navigation, or direct alias execution.",
	RunE:  runRun,
}

func init() {
	runCmd.Flags().String("preset", "", "Preset name for parameterized commands")
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	configDir, err := config.Dir()
	if err != nil {
		return err
	}

	globalPath := filepath.Join(configDir, "commands.yaml")
	projectPath := ".bro.yaml"

	shortcuts, err := runner.LoadShortcuts(globalPath, projectPath)
	if err != nil {
		return fmt.Errorf("load shortcuts: %w", err)
	}

	if len(shortcuts) == 0 {
		fmt.Println("No shortcuts configured.")
		fmt.Printf("Add shortcuts to %s or .bro.yaml in your project root.\n", globalPath)
		fmt.Println("\nExample commands.yaml:")
		fmt.Println("  shortcuts:")
		fmt.Println("    - name: test")
		fmt.Println("      command: go test ./...")
		fmt.Println("      description: Run all tests")
		return nil
	}

	// Alias mode: bro run <name>
	if len(args) > 0 {
		name := args[0]
		for _, s := range shortcuts {
			if s.Name == name {
				command := s.Command

				// Apply preset if specified
				presetName, _ := cmd.Flags().GetString("preset")
				if presetName != "" && s.Presets != nil {
					preset, ok := s.Presets[presetName]
					if !ok {
						available := make([]string, 0, len(s.Presets))
						for k := range s.Presets {
							available = append(available, k)
						}
						return fmt.Errorf("preset %q not found. Available: %v", presetName, available)
					}
					command = runner.ApplyPreset(command, preset)
				}

				fmt.Printf("Running: %s\n", ui.DimStyle.Render(command))
				return runner.Execute(command)
			}
		}
		return fmt.Errorf("shortcut %q not found. Run 'bro run' to see available shortcuts", name)
	}

	// Interactive mode: fuzzy picker
	items := make([]ui.PickerItem, len(shortcuts))
	for i, s := range shortcuts {
		// Serialize presets info for display if present
		desc := s.Description
		if len(s.Presets) > 0 {
			presetNames := make([]string, 0, len(s.Presets))
			for k := range s.Presets {
				presetNames = append(presetNames, k)
			}
			data, _ := json.Marshal(presetNames)
			desc += fmt.Sprintf(" (presets: %s)", string(data))
		}
		items[i] = ui.PickerItem{
			Name:        s.Name,
			Description: desc,
			Value:       s.Command,
		}
	}

	selected, err := ui.RunPicker(items)
	if err != nil {
		return fmt.Errorf("picker: %w", err)
	}

	if selected == nil {
		return nil
	}

	fmt.Printf("Running: %s\n", ui.DimStyle.Render(selected.Value))
	return runner.Execute(selected.Value)
}
