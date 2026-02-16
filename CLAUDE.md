# bro — Developer Productivity CLI

Go CLI that reduces daily dev friction: smart commit messages, command shortcuts, standup summaries, and AI-assisted Q&A.

## Tech Stack

- Language: Go 1.22+
- CLI framework: cobra + viper (config)
- Terminal UI: bubbletea + lipgloss (interactive prompts, fuzzy picker)
- LLM: Anthropic (default), OpenAI (alternative) — direct HTTP via `net/http`, no SDKs
- Git: shell out to `git` via `os/exec` (not go-git)
- Config format: YAML
- Module path: `github.com/<username>/bro-cli`

## Commands (V1)

- `bro commit` — generate commit message from `git diff --staged`, interactive accept/edit/regenerate flow
- `bro run` — fuzzy-searchable command shortcuts (arrow key navigation) + direct alias mode
- `bro standup` — scan git logs across configured repos, LLM-summarize, save as markdown
- `bro ask` — general LLM queries, supports stdin pipe and `-f` file flag, streams response
- `bro config` — manage config, profiles, API keys

## Project Structure

```
bro-cli/
├── cmd/                    # Cobra command definitions
│   ├── root.go             # Root command, setup wizard on first run
│   ├── commit.go
│   ├── run.go
│   ├── standup.go
│   ├── ask.go
│   └── config.go
├── internal/
│   ├── llm/                # LLM provider abstraction
│   │   ├── provider.go     # Provider interface: Generate(prompt, systemPrompt) (string, error)
│   │   ├── anthropic.go    # Anthropic Messages API implementation
│   │   └── openai.go       # OpenAI Chat Completions implementation
│   ├── git/                # Git operations via os/exec
│   │   ├── diff.go         # git diff --staged
│   │   ├── log.go          # git log parsing for standup
│   │   └── branch.go       # Branch name parsing with regex
│   ├── config/             # Config loading, profiles, directory matching
│   │   ├── config.go       # Viper-based config struct + loading
│   │   └── profile.go      # Profile matching via directory glob patterns
│   ├── runner/             # Command shortcuts
│   │   ├── shortcuts.go    # Load from global + project-level YAML
│   │   └── executor.go     # Template params, presets, shell execution
│   ├── standup/            # Standup generation
│   │   └── generator.go    # Repo scanning, log collection, LLM summary
│   └── ui/                 # Terminal UI components (bubbletea)
│       ├── picker.go       # Scrollable fuzzy list for `bro run`
│       ├── commit_flow.go  # Accept/edit/regenerate prompt for `bro commit`
│       └── styles.go       # Shared lipgloss styles
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Config Layout

User config lives in `~/.bro/`. Project-level overrides in `.bro.yaml` at repo root.

```
~/.bro/
├── config.yaml        # Main config (LLM keys, profiles, project mappings)
├── commands.yaml      # Global command shortcuts
├── standups/          # Generated standup markdown files (YYYY-MM-DD.md)
└── templates/         # Custom commit message templates
```

## Key Design Decisions

- **Profile system**: Profiles define commit behavior (template vs conventional commits). Mapped to directories via glob patterns in config (e.g., `~/work/*` → org profile, `~/personal/*` → personal profile).
- **Org commit template**: `{project}-{ticket} #comment {message}` — ticket/project extracted from branch name via configurable regex `branch_pattern`.
- **Personal commit style**: Conventional commits (feat/fix/chore/docs/refactor/test/ci). LLM picks the type from the diff.
- **LLM provider interface**: Keep it simple — `Generate(ctx, prompt, systemPrompt string) (string, error)` and `StreamGenerate(ctx, prompt, systemPrompt string) (<-chan string, error)`. Provider selected via config, API keys from config or env vars (`BRO_ANTHROPIC_API_KEY`, `BRO_OPENAI_API_KEY`).
- **`bro run` interactive mode**: Bubbletea-powered list with arrow key scrolling, type-to-filter, enter to execute. Loads shortcuts from both `~/.bro/commands.yaml` (global) and `.bro.yaml` (project-level, takes precedence).
- **`bro run` alias mode**: `bro run <name>` executes directly. Supports `--preset` flag for parameterized commands.
- **`bro standup`**: Semi-automatic (user triggers it). Scans configured repos for git activity since last working day (auto-detects weekends). Saves to `~/.bro/standups/YYYY-MM-DD.md`.
- **`bro ask`**: Reads stdin if piped, `-f` flag for file context. Streams LLM response token-by-token to terminal.
- **No browser/network features**: This is a local-first CLI. Only outbound calls are to LLM APIs.

## Code Style

- Use standard Go conventions (gofmt, go vet)
- Error handling: wrap errors with `fmt.Errorf("context: %w", err)`, don't swallow errors silently
- Use `internal/` to keep packages unexported
- Prefer returning errors over panicking
- Keep cobra command files thin — business logic goes in `internal/` packages
- Use table-driven tests

## Commands Reference

```bash
go run .                    # Run locally
go build -o bro .           # Build binary
go test ./...               # Run all tests
go vet ./...                # Lint
```

## Important

- Anthropic is the DEFAULT LLM provider. Use the Messages API (`/v1/messages`), not the legacy completions endpoint.
- For OpenAI, use Chat Completions API (`/v1/chat/completions`).
- Both providers: send streaming requests for `bro ask`, non-streaming for `bro commit` and `bro standup`.
- The interactive UI (bubbletea) for `bro run` must support: arrow up/down to scroll, typing to filter, enter to select and execute, escape/q to quit.
- The interactive UI for `bro commit` must support: enter to accept, `e` to open in `$EDITOR`, `r` to regenerate, `q` to quit.
- When no staged changes exist for `bro commit`, show a helpful error suggesting `git add`.
- When LLM API is unavailable, `bro commit` should fall back to template-only mode where user fills placeholders manually.
- First run with no config should launch an interactive setup wizard (choose provider, enter API key, set default profile).