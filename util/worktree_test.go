package util

import "testing"

func TestSanitizeForSessionName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "branch with single slash",
			input:    "feature/ABC123",
			expected: "feature-ABC123",
		},
		{
			name:     "branch with multiple slashes",
			input:    "feature/api/user-endpoint",
			expected: "feature-api-user-endpoint",
		},
		{
			name:     "branch with colon",
			input:    "feature/bug-fix:123",
			expected: "feature-bug-fix-123",
		},
		{
			name:     "branch with dots",
			input:    "release/v1.2.3",
			expected: "release-v1_2_3",
		},
		{
			name:     "branch with spaces",
			input:    "feature/my branch",
			expected: "feature-my_branch",
		},
		{
			name:     "branch with mixed special characters",
			input:    "feature/bug:fix.v2.0",
			expected: "feature-bug-fix_v2_0",
		},
		{
			name:     "simple branch name",
			input:    "main",
			expected: "main",
		},
		{
			name:     "branch with dashes only",
			input:    "feature-branch",
			expected: "feature-branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeForSessionName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeForSessionName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
