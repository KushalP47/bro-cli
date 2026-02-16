package git

import "testing"

func TestParseBranch(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		pattern string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "jira-style branch",
			branch:  "feature/PROJ-123-add-login",
			pattern: `(?P<type>\w+)/(?P<project>[A-Z]+)-(?P<ticket>\d+)-(?P<desc>.+)`,
			want:    map[string]string{"type": "feature", "project": "PROJ", "ticket": "123", "desc": "add-login"},
		},
		{
			name:    "simple ticket branch",
			branch:  "PROJ-456-fix-bug",
			pattern: `(?P<project>[A-Z]+)-(?P<ticket>\d+)`,
			want:    map[string]string{"project": "PROJ", "ticket": "456"},
		},
		{
			name:    "no match",
			branch:  "main",
			pattern: `(?P<project>[A-Z]+)-(?P<ticket>\d+)`,
			wantErr: true,
		},
		{
			name:    "invalid regex",
			branch:  "test",
			pattern: `(?P<broken`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBranch(tt.branch, tt.pattern)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for key, want := range tt.want {
				got, ok := result[key]
				if !ok {
					t.Errorf("missing key %q", key)
					continue
				}
				if got != want {
					t.Errorf("key %q: expected %q, got %q", key, want, got)
				}
			}
		})
	}
}
