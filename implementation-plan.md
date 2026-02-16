bro-cli V1 Implementation Plan                                                                                                                       │
│                                                                                                                                                      │
│ Context                                                                                                                                              │
│                                                                                                                                                      │
│ The bro-cli project is scaffolded with cobra command stubs, working git utilities, config loading, and a defined LLM provider interface. All five    │
│ commands (commit, run, standup, ask, config) print "not yet implemented". The LLM providers, UI components, shortcut loading, standup generation,    │
│ and setup wizard all need to be built. This plan implements the complete V1 feature set.                                                             │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 0: Add Dependencies                                                                                                                            │
│                                                                                                                                                      │
│ Install required Go packages not yet in go.mod:                                                                                                      │
│                                                                                                                                                      │
│ go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/charmbracelet/bubbles                                         │
│ go mod tidy                                                                                                                                          │
│                                                                                                                                                      │
│ Verify: go build ./... compiles.                                                                                                                     │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 1: LLM Providers                                                                                                                               │
│                                                                                                                                                      │
│ All LLM-dependent commands (commit, ask, standup) depend on this.                                                                                    │
│                                                                                                                                                      │
│ 1A: Anthropic Provider — internal/llm/anthropic.go                                                                                                   │
│                                                                                                                                                      │
│ - Generate(): POST to https://api.anthropic.com/v1/messages with x-api-key, anthropic-version: 2023-06-01. Parse response.content[0].text. Max       │
│ tokens: 1024.                                                                                                                                        │
│ - StreamGenerate(): Same endpoint with "stream": true. Parse SSE events — extract text from content_block_delta events. Return channel, read in      │
│ goroutine. Close on message_stop.                                                                                                                    │
│                                                                                                                                                      │
│ 1B: OpenAI Provider — internal/llm/openai.go                                                                                                         │
│                                                                                                                                                      │
│ - Generate(): POST to https://api.openai.com/v1/chat/completions with Bearer token. System message + user message. Parse choices[0].message.content. │
│ - StreamGenerate(): Same with "stream": true. Parse SSE data: lines, extract choices[0].delta.content. Close on data: [DONE].                        │
│                                                                                                                                                      │
│ 1C: Provider Factory — internal/llm/provider.go                                                                                                      │
│                                                                                                                                                      │
│ - Add NewProvider(cfg *config.Config) (Provider, error) — selects provider by cfg.Provider, resolves API key from config or env vars                 │
│ (BRO_ANTHROPIC_API_KEY, BRO_OPENAI_API_KEY), returns configured provider.                                                                            │
│ - Update NewAnthropicProvider/NewOpenAIProvider to accept (apiKey, model string).                                                                    │
│                                                                                                                                                      │
│ Verify: go build ./... compiles. Manual test with a real API key.                                                                                    │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 2: UI Components (bubbletea + lipgloss)                                                                                                        │
│                                                                                                                                                      │
│ 2A: Shared Styles — internal/ui/styles.go                                                                                                            │
│                                                                                                                                                      │
│ - Define lipgloss styles: TitleStyle, SelectedStyle, NormalStyle, DimStyle, SuccessStyle, ErrorStyle, HelpStyle.                                     │
│                                                                                                                                                      │
│ 2B: Commit Flow — internal/ui/commit_flow.go                                                                                                         │
│                                                                                                                                                      │
│ - Bubbletea model showing the generated message + key options                                                                                        │
│ - Keys: enter → accept, e → edit, r → regenerate, q/esc → quit                                                                                       │
│ - Export RunCommitFlow(message string) (choice string, err error)                                                                                    │
│                                                                                                                                                      │
│ 2C: Fuzzy Picker — internal/ui/picker.go                                                                                                             │
│                                                                                                                                                      │
│ - Bubbletea model with type-to-filter, arrow key navigation, enter to select                                                                         │
│ - PickerItem{Name, Description, Value} struct                                                                                                        │
│ - Case-insensitive substring filtering on Name + Description                                                                                         │
│ - Export RunPicker(items []PickerItem) (*PickerItem, error)                                                                                          │
│                                                                                                                                                      │
│ Verify: go build ./... compiles.                                                                                                                     │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 3: bro ask — Simplest LLM Command                                                                                                              │
│                                                                                                                                                      │
│ File: cmd/ask.go                                                                                                                                     │
│                                                                                                                                                      │
│ 1. Load config → create provider via llm.NewProvider                                                                                                 │
│ 2. Gather question from args, stdin (if piped), and -f file flag                                                                                     │
│ 3. System prompt: "You are a helpful developer assistant. Be concise and precise."                                                                   │
│ 4. Call provider.StreamGenerate() — print each chunk to stdout                                                                                       │
│ 5. Handle errors: no API key → suggest bro config init or env var                                                                                    │
│                                                                                                                                                      │
│ Verify: go run . ask "what is Go?" streams response. echo "error text" | go run . ask "why?" works. go run . ask -f main.go "explain this" works.    │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 4: bro commit — Smart Commit Messages                                                                                                          │
│                                                                                                                                                      │
│ File: cmd/commit.go                                                                                                                                  │
│                                                                                                                                                      │
│ 1. git.StagedDiff() — if empty, error with "No staged changes. Use 'git add' to stage files first."                                                  │
│ 2. Load config → create provider (if fails, enter fallback mode)                                                                                     │
│ 3. Match profile via config.MatchProfile(cfg.Profiles, cwd)                                                                                          │
│ 4. Build LLM prompt based on profile:                                                                                                                │
│   - Conventional: system prompt instructs conventional format, user prompt includes diff                                                             │
│   - Template: extract ticket/project from branch via git.ParseBranch(), LLM generates message body, apply template                                   │
│ 5. Interactive loop using ui.RunCommitFlow(message):                                                                                                 │
│   - accept → git commit -m "<message>"                                                                                                               │
│   - edit → open $EDITOR with temp file, read back, commit                                                                                            │
│   - regenerate → call LLM again, loop                                                                                                                │
│   - quit → exit                                                                                                                                      │
│ 6. Fallback mode (no LLM): open template in $EDITOR for manual fill                                                                                  │
│ 7. Truncate large diffs (~8000 chars) to fit LLM context                                                                                             │
│                                                                                                                                                      │
│ Verify: Stage changes, run go run . commit, test all four choices.                                                                                   │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 5: bro run — Command Shortcuts                                                                                                                 │
│                                                                                                                                                      │
│ 5A: Shortcuts Loading — internal/runner/shortcuts.go                                                                                                 │
│                                                                                                                                                      │
│ - Implement LoadShortcuts(globalPath, projectPath string) ([]Shortcut, error)                                                                        │
│ - Parse YAML from both paths, merge (project overrides global by name)                                                                               │
│ - Add Presets map[string]map[string]string to Shortcut struct                                                                                        │
│ - Add ApplyPreset(command string, preset map[string]string) string for {{placeholder}} substitution                                                  │
│                                                                                                                                                      │
│ 5B: Command — cmd/run.go                                                                                                                             │
│                                                                                                                                                      │
│ - Load shortcuts from ~/.bro/commands.yaml + .bro.yaml                                                                                               │
│ - Alias mode (bro run <name>): find by name, apply --preset if set, execute                                                                          │
│ - Interactive mode (bro run): convert to PickerItems, run ui.RunPicker(), execute selection                                                          │
│ - No shortcuts found → helpful message                                                                                                               │
│                                                                                                                                                      │
│ Verify: Create ~/.bro/commands.yaml with test commands. Test both modes.                                                                             │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 6: bro standup — Daily Summary                                                                                                                 │
│                                                                                                                                                      │
│ 6A: Generator — internal/standup/generator.go                                                                                                        │
│                                                                                                                                                      │
│ - Update Generator to accept llm.Provider                                                                                                            │
│ - Determine "since" date: Monday → last Friday, otherwise yesterday                                                                                  │
│ - Loop through repos, call git.Log(), aggregate                                                                                                      │
│ - Send to LLM with standup summarization system prompt (non-streaming)                                                                               │
│ - Return summary string                                                                                                                              │
│                                                                                                                                                      │
│ 6B: Command — cmd/standup.go                                                                                                                         │
│                                                                                                                                                      │
│ - Load config, create provider                                                                                                                       │
│ - Check cfg.StandupRepos is non-empty                                                                                                                │
│ - Call generator.Generate(), print preview                                                                                                           │
│ - Save to ~/.bro/standups/YYYY-MM-DD.md                                                                                                              │
│                                                                                                                                                      │
│ Verify: Configure repos in ~/.bro/config.yaml, run go run . standup.                                                                                 │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 7: bro config + Setup Wizard                                                                                                                   │
│                                                                                                                                                      │
│ 7A: Config Subcommands — cmd/config.go                                                                                                               │
│                                                                                                                                                      │
│ - bro config show — print current config as YAML                                                                                                     │
│ - bro config set <key> <value> — update config via viper                                                                                             │
│ - bro config init — launch setup wizard                                                                                                              │
│                                                                                                                                                      │
│ 7B: Setup Wizard — triggered in cmd/root.go                                                                                                          │
│                                                                                                                                                      │
│ - Detect first run (no ~/.bro/config.yaml)                                                                                                           │
│ - Interactive prompts: choose provider → enter API key → set model                                                                                   │
│ - Write config via viper                                                                                                                             │
│ - Only trigger for commands needing LLM, not for config init itself                                                                                  │
│                                                                                                                                                      │
│ Verify: Delete config, run go run . ask "hello" → wizard launches. go run . config show displays config.                                             │
│                                                                                                                                                      │
│ ---                                                                                                                                                  │
│ Phase 8: Testing & Polish                                                                                                                            │
│                                                                                                                                                      │
│ Tests (table-driven):                                                                                                                                │
│                                                                                                                                                      │
│ - internal/llm/anthropic_test.go — mock HTTP server, test request/response parsing                                                                   │
│ - internal/llm/openai_test.go — same                                                                                                                 │
│ - internal/llm/provider_test.go — test factory with various configs                                                                                  │
│ - internal/runner/shortcuts_test.go — test loading, merging, template substitution                                                                   │
│ - internal/standup/generator_test.go — test date logic, mock provider                                                                                │
│ - internal/git/branch_test.go — test ParseBranch patterns                                                                                            │
│ - internal/config/profile_test.go — test MatchProfile with globs                                                                                     │
│                                                                                                                                                      │
│ Edge Cases:                                                                                                                                          │
│                                                                                                                                                      │
│ - LLM timeout (30s context for Generate, 60s for StreamGenerate)                                                                                     │
│ - Invalid API key → clear error message                                                                                                              │
│ - Network errors → wrapped with context                                                                                                              │
│ - Unreachable standup repos → skip with warning, don't crash                                                                                         │
│                                                                                                                                                      │
│ Polish:                                                                                                                                              │
│                                                                                                                                                      │
│ - Graceful ctrl+c handling (context cancellation)                                                                                                    │
│ - bro version command with build-time version via ldflags                                                                                            │
│                                                                                                                                                      │
│ Verify: go test ./... passes. go vet ./... clean. Manual end-to-end test of all 5 commands.    