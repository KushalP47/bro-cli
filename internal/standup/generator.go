package standup

// Generator handles scanning repos and producing standup summaries.
type Generator struct {
	Repos []string
}

// NewGenerator creates a standup generator for the given repo paths.
func NewGenerator(repos []string) *Generator {
	return &Generator{Repos: repos}
}

// Generate scans repos for recent git activity and produces a summary.
func (g *Generator) Generate() (string, error) {
	// TODO: scan repos, collect logs, summarize with LLM
	return "", nil
}
