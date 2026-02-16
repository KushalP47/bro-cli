package config

import "testing"

func TestMatchProfile(t *testing.T) {
	profiles := []Profile{
		{Name: "work", DirPattern: "/Users/dev/work/*", CommitStyle: "template"},
		{Name: "personal", DirPattern: "/Users/dev/personal/*", CommitStyle: "conventional"},
	}

	tests := []struct {
		name     string
		dir      string
		wantName string
		wantNil  bool
	}{
		{
			name:     "match work profile",
			dir:      "/Users/dev/work/myproject",
			wantName: "work",
		},
		{
			name:     "match personal profile",
			dir:      "/Users/dev/personal/sideproject",
			wantName: "personal",
		},
		{
			name:    "no match",
			dir:     "/Users/dev/other/something",
			wantNil: true,
		},
		{
			name:    "empty profiles",
			dir:     "/some/path",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p []Profile
			if tt.name != "empty profiles" {
				p = profiles
			}

			result := MatchProfile(p, tt.dir)

			if tt.wantNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("expected non-nil profile")
			}

			if result.Name != tt.wantName {
				t.Errorf("expected profile %q, got %q", tt.wantName, result.Name)
			}
		})
	}
}
