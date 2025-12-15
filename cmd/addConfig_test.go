package cmd

import (
	"testing"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBuildConfigWorktreeDir(t *testing.T) {
	tests := []struct {
		name                    string
		homePath                string
		configWorktreeTargetDir string
		branchName              string
		expected                string
		description             string
	}{
		{
			name:                    "both paths end with slash",
			homePath:                "/home/user/",
			configWorktreeTargetDir: "worktrees/",
			branchName:              "feature-branch",
			expected:                "/home/user/worktrees/feature-branch",
			description:             "homePath ends with '/' and configWorktreeTargetDir ends with '/'",
		},
		{
			name:                    "homePath ends with slash, configWorktreeTargetDir does not",
			homePath:                "/home/user/",
			configWorktreeTargetDir: "worktrees",
			branchName:              "feature-branch",
			expected:                "/home/user/worktrees/feature-branch",
			description:             "homePath ends with '/' and configWorktreeTargetDir doesn't end with '/'",
		},
		{
			name:                    "homePath does not end with slash, configWorktreeTargetDir ends with slash",
			homePath:                "/home/user",
			configWorktreeTargetDir: "worktrees/",
			branchName:              "feature-branch",
			expected:                "/home/user/worktrees/feature-branch",
			description:             "homePath doesn't end with '/' and configWorktreeTargetDir ends with '/'",
		},
		{
			name:                    "neither path ends with slash",
			homePath:                "/home/user",
			configWorktreeTargetDir: "worktrees",
			branchName:              "feature-branch",
			expected:                "/home/user/worktrees/feature-branch",
			description:             "homePath doesn't end with '/' and configWorktreeTargetDir doesn't end with '/'",
		},
		{
			name:                    "empty configWorktreeTargetDir with trailing slash on homePath",
			homePath:                "/home/user/",
			configWorktreeTargetDir: "",
			branchName:              "feature-branch",
			expected:                "/home/user/feature-branch",
			description:             "configWorktreeTargetDir is empty, homePath ends with '/'",
		},
		{
			name:                    "empty configWorktreeTargetDir without trailing slash on homePath",
			homePath:                "/home/user",
			configWorktreeTargetDir: "",
			branchName:              "feature-branch",
			expected:                "/home/user/feature-branch",
			description:             "configWorktreeTargetDir is empty, homePath doesn't end with '/'",
		},
		{
			name:                    "multiple slashes in paths",
			homePath:                "/home/user//",
			configWorktreeTargetDir: "//worktrees//",
			branchName:              "feature-branch",
			expected:                "/home/user/worktrees/feature-branch",
			description:             "paths with multiple slashes (should be cleaned up by filepath.Join)",
		},
		{
			name:                    "branch name with special characters",
			homePath:                "/home/user",
			configWorktreeTargetDir: "worktrees",
			branchName:              "feature/special-branch_123",
			expected:                "/home/user/worktrees/feature/special-branch_123",
			description:             "branch name containing special characters and slashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildConfigWorktreeDir(tt.homePath, tt.configWorktreeTargetDir, tt.branchName)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestGetWorktreeName(t *testing.T) {
	tests := []struct {
		name        string
		config      *com.AddConfig
		expected    string
		description string
	}{
		{
			name: "uses specified worktree name when provided",
			config: &com.AddConfig{
				Flags: com.AddCmdFlags{
					SpecifiedWorktreeName: stringPtr("custom-name"),
				},
				GitInfo: com.GitInfo{
					NewBranchName: "feature-branch",
				},
			},
			expected:    "custom-name",
			description: "should use the specified worktree name when provided",
		},
		{
			name: "uses branch name when no specified name",
			config: &com.AddConfig{
				Flags: com.AddCmdFlags{
					SpecifiedWorktreeName: nil,
				},
				GitInfo: com.GitInfo{
					NewBranchName: "feature-branch",
				},
			},
			expected:    "feature-branch",
			description: "should use branch name when no specific name is provided",
		},
		{
			name: "uses branch name when specified name is empty",
			config: &com.AddConfig{
				Flags: com.AddCmdFlags{
					SpecifiedWorktreeName: stringPtr(""),
				},
				GitInfo: com.GitInfo{
					NewBranchName: "feature-branch",
				},
			},
			expected:    "feature-branch",
			description: "should use branch name when specified name is empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getWorktreeName(tt.config)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// Helper function to create string pointers for tests
func stringPtr(s string) *string {
	return &s
}

func TestFindRepoByBareRepoName(t *testing.T) {
	tests := []struct {
		name         string
		bareRepoName string
		configSetup  func()
		expected     string
		description  string
	}{
		{
			name:         "finds repo with matching bareRepoName",
			bareRepoName: ".bare",
			configSetup: func() {
				viper.Reset()
				viper.Set("repos.myproject.bareRepoName", ".bare")
				viper.Set("repos.myproject.defaultBranch", "main")
			},
			expected:    "repos.myproject",
			description: "should find repo when bareRepoName matches",
		},
		{
			name:         "returns empty when no match found",
			bareRepoName: ".bare",
			configSetup: func() {
				viper.Reset()
				viper.Set("repos.myproject.bareRepoName", "_bare")
				viper.Set("repos.myproject.defaultBranch", "main")
			},
			expected:    "",
			description: "should return empty string when no bareRepoName matches",
		},
		{
			name:         "finds correct repo among multiple repos",
			bareRepoName: ".bare",
			configSetup: func() {
				viper.Reset()
				viper.Set("repos.project1.bareRepoName", "_bare")
				viper.Set("repos.project1.defaultBranch", "main")
				viper.Set("repos.project2.bareRepoName", ".bare")
				viper.Set("repos.project2.defaultBranch", "develop")
				viper.Set("repos.project3.bareRepoName", "-bare")
				viper.Set("repos.project3.defaultBranch", "master")
			},
			expected:    "repos.project2",
			description: "should find the correct repo when multiple repos exist",
		},
		{
			name:         "returns empty when repo has no bareRepoName configured",
			bareRepoName: ".bare",
			configSetup: func() {
				viper.Reset()
				viper.Set("repos.myproject.defaultBranch", "main")
			},
			expected:    "",
			description: "should return empty when repo doesn't have bareRepoName configured",
		},
		{
			name:         "handles underscore bare naming convention",
			bareRepoName: "myproject_bare",
			configSetup: func() {
				viper.Reset()
				viper.Set("repos.myproject.bareRepoName", "myproject_bare")
				viper.Set("repos.myproject.defaultBranch", "main")
			},
			expected:    "repos.myproject",
			description: "should work with traditional _bare naming convention",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup config for this test
			tt.configSetup()

			result := findRepoByBareRepoName(tt.bareRepoName)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}

	// Clean up
	viper.Reset()
}
