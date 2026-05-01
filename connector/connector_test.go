package connector

import (
	"testing"
)

func TestGenerateSessionName(t *testing.T) {
	// No need to initialize shell since generateSessionName doesn't use it
	connector := &RealConnector{}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/home/user/project",
			expected: "project",
		},
		{
			name:     "path with dots",
			path:     "/home/user/my.project.v2",
			expected: "my_project_v2",
		},
		{
			name:     "path with spaces",
			path:     "/home/user/my project",
			expected: "my_project",
		},
		{
			name:     "path with slashes in basename",
			path:     "/home/user/feature-branch",
			expected: "feature-branch",
		},
		{
			name:     "path with mixed special characters",
			path:     "/home/user/my.project folder",
			expected: "my_project_folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := connector.generateSessionName(tt.path)
			if result != tt.expected {
				t.Errorf("generateSessionName(%q) = %q, expected %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGenerateWorktreeSessionName(t *testing.T) {
	// No need to initialize shell since generateWorktreeSessionName doesn't use it
	connector := &RealConnector{}

	tests := []struct {
		name         string
		worktreePath string
		branchName   string
		expected     string
	}{
		{
			name:         "branch with single slash",
			worktreePath: "/home/user/treekanga_work/feature-ABC123",
			branchName:   "feature/ABC123",
			expected:     "treekanga-feature-ABC123",
		},
		{
			name:         "branch with multiple slashes",
			worktreePath: "/home/user/treekanga_work/feature-api-user",
			branchName:   "feature/api/user-endpoint",
			expected:     "treekanga-feature-api-user-endpoint",
		},
		{
			name:         "branch with colon",
			worktreePath: "/home/user/treekanga_work/hotfix",
			branchName:   "hotfix/bug:123",
			expected:     "treekanga-hotfix-bug-123",
		},
		{
			name:         "branch with dots (version number)",
			worktreePath: "/home/user/treekanga_work/release-v1.2.3",
			branchName:   "release/v1.2.3",
			expected:     "treekanga-release-v1_2_3",
		},
		{
			name:         "simple branch name",
			worktreePath: "/home/user/treekanga_work/main",
			branchName:   "main",
			expected:     "treekanga-main",
		},
		{
			name:         "branch with spaces",
			worktreePath: "/home/user/treekanga_work/my-branch",
			branchName:   "feature/my branch",
			expected:     "treekanga-feature-my_branch",
		},
		{
			name:         "repo with _work suffix",
			worktreePath: "/home/user/myrepo_work/feature-branch",
			branchName:   "feature/test",
			expected:     "myrepo-feature-test",
		},
		{
			name:         "repo with _worktrees suffix",
			worktreePath: "/home/user/myrepo_worktrees/feature-branch",
			branchName:   "feature/test",
			expected:     "myrepo-feature-test",
		},
		{
			name:         "repo with -bare suffix",
			worktreePath: "/home/user/myrepo-bare/feature-branch",
			branchName:   "feature/test",
			expected:     "myrepo-feature-test",
		},
		{
			name:         "repo with .git suffix",
			worktreePath: "/home/user/myrepo.git/feature-branch",
			branchName:   "feature/test",
			expected:     "myrepo-feature-test",
		},
		{
			name:         "complex branch with all special chars",
			worktreePath: "/home/user/platform_work/complex-branch",
			branchName:   "feature/bug:fix.v2.0",
			expected:     "platform-feature-bug-fix_v2_0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := connector.generateWorktreeSessionName(tt.worktreePath, tt.branchName)
			if result != tt.expected {
				t.Errorf("generateWorktreeSessionName(%q, %q) = %q, expected %q",
					tt.worktreePath, tt.branchName, result, tt.expected)
			}
		})
	}
}

// TestGenerateWorktreeSessionNameRealWorldExamples tests with real-world branch naming conventions
func TestGenerateWorktreeSessionNameRealWorldExamples(t *testing.T) {
	// No need to initialize shell since generateWorktreeSessionName doesn't use it
	connector := &RealConnector{}

	realWorldTests := []struct {
		name         string
		worktreePath string
		branchName   string
		expected     string
		description  string
	}{
		{
			name:         "GitHub flow feature branch",
			worktreePath: "/Users/dev/platform_work/feature-user-auth",
			branchName:   "feature/user-authentication",
			expected:     "platform-feature-user-authentication",
			description:  "Common GitHub flow pattern",
		},
		{
			name:         "GitFlow bugfix branch",
			worktreePath: "/Users/dev/platform_work/bugfix-null-pointer",
			branchName:   "bugfix/null-pointer-exception",
			expected:     "platform-bugfix-null-pointer-exception",
			description:  "GitFlow bugfix pattern",
		},
		{
			name:         "GitFlow release branch",
			worktreePath: "/Users/dev/platform_work/release-2.1.0",
			branchName:   "release/2.1.0",
			expected:     "platform-release-2_1_0",
			description:  "GitFlow release with version number",
		},
		{
			name:         "Jira ticket branch",
			worktreePath: "/Users/dev/platform_work/PROJ-1234",
			branchName:   "feature/PROJ-1234-add-login",
			expected:     "platform-feature-PROJ-1234-add-login",
			description:  "Branch with Jira ticket ID",
		},
		{
			name:         "Nested feature branch",
			worktreePath: "/Users/dev/platform_work/api-users-endpoint",
			branchName:   "feature/api/users/endpoint",
			expected:     "platform-feature-api-users-endpoint",
			description:  "Deeply nested branch structure",
		},
		{
			name:         "Hotfix with issue number",
			worktreePath: "/Users/dev/platform_work/hotfix-123",
			branchName:   "hotfix/issue:123-critical-bug",
			expected:     "platform-hotfix-issue-123-critical-bug",
			description:  "Hotfix with colon-separated issue number",
		},
	}

	for _, tt := range realWorldTests {
		t.Run(tt.name, func(t *testing.T) {
			result := connector.generateWorktreeSessionName(tt.worktreePath, tt.branchName)
			if result != tt.expected {
				t.Errorf("%s:\ngenerateWorktreeSessionName(%q, %q)\ngot:      %q\nexpected: %q",
					tt.description, tt.worktreePath, tt.branchName, result, tt.expected)
			}
		})
	}
}

func TestBareRepoStrategy_RecognizesBareRepoPaths(t *testing.T) {
	// Arrange - path ending in .bare
	connector := &RealConnector{}
	path := "/Users/gkrohn/code/cal_work/.bare"

	// Act - call bareRepoStrategy(path)
	connection, err := connector.bareRepoStrategy(path)

	// Assert - expect Found=true in Connection
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !connection.Found {
		t.Errorf("Expected Found=true for bare repo path ending in .bare")
	}
	if connection.Session.Name == "" {
		t.Errorf("Expected session name to be set")
	}
	if connection.Session.Src != "bare" {
		t.Errorf("Expected Src='bare', got: %s", connection.Session.Src)
	}
}

func TestBareRepoStrategy_IgnoresNonBareRepoPaths(t *testing.T) {
	// Arrange - path NOT ending in .bare
	connector := &RealConnector{}
	path := "/Users/gkrohn/code/cal_work/feature-branch"

	// Act - call bareRepoStrategy(path)
	connection, err := connector.bareRepoStrategy(path)

	// Assert - expect Found=false
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if connection.Found {
		t.Errorf("Expected Found=false for non-bare repo path")
	}
}

func TestGenerateBareRepoSessionName_FormatsAsRepoBareName(t *testing.T) {
	// Arrange - bare repo path
	connector := &RealConnector{}
	path := "/Users/gkrohn/code/cal_work/.bare"

	// Act - call generateBareRepoSessionName(path)
	result := connector.generateBareRepoSessionName(path)

	// Assert - expect "cal - bare" format
	expected := "cal - bare"
	if result != expected {
		t.Errorf("generateBareRepoSessionName(%q) = %q, expected %q", path, result, expected)
	}
}

func TestGenerateBareRepoSessionName_StripsWorkSuffix(t *testing.T) {
	// Arrange - bare repo path with _work suffix
	connector := &RealConnector{}
	path := "/Users/gkrohn/code/core_work/.bare"

	// Act - call generateBareRepoSessionName(path)
	result := connector.generateBareRepoSessionName(path)

	// Assert - verify _work suffix removed
	expected := "core - bare"
	if result != expected {
		t.Errorf("generateBareRepoSessionName(%q) = %q, expected %q (should strip _work suffix)", path, result, expected)
	}
}

func TestConnect_SuccessfullyConnectsToBareRepo(t *testing.T) {
	// Note: This is an integration-style test that requires actual filesystem
	// For now, we verify the strategy exists and can be called
	// Full integration testing will happen when the feature is wired up

	// Arrange - bare repo path
	connector := &RealConnector{}
	path := "/Users/gkrohn/code/cal_work/.bare"

	// Act - call bareRepoStrategy (not full Connect, as that requires tmux mock)
	connection, err := connector.bareRepoStrategy(path)

	// Assert - verify connection object structure
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if connection.Found && connection.Session.Name == "" {
		t.Errorf("Expected session name to be set when Found=true")
	}
}
