package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConfigWorktreeDir(t *testing.T) {
	tests := []struct {
		name                     string
		homePath                 string
		configWorktreeTargetDir  string
		branchName              string
		expected                string
		description             string
	}{
		{
			name:                    "both paths end with slash",
			homePath:                "/home/user/",
			configWorktreeTargetDir: "worktrees/",
			branchName:             "feature-branch",
			expected:               "/home/user/worktrees/feature-branch",
			description:            "homePath ends with '/' and configWorktreeTargetDir ends with '/'",
		},
		{
			name:                    "homePath ends with slash, configWorktreeTargetDir does not",
			homePath:                "/home/user/",
			configWorktreeTargetDir: "worktrees",
			branchName:             "feature-branch",
			expected:               "/home/user/worktrees/feature-branch",
			description:            "homePath ends with '/' and configWorktreeTargetDir doesn't end with '/'",
		},
		{
			name:                    "homePath does not end with slash, configWorktreeTargetDir ends with slash",
			homePath:                "/home/user",
			configWorktreeTargetDir: "worktrees/",
			branchName:             "feature-branch",
			expected:               "/home/user/worktrees/feature-branch",
			description:            "homePath doesn't end with '/' and configWorktreeTargetDir ends with '/'",
		},
		{
			name:                    "neither path ends with slash",
			homePath:                "/home/user",
			configWorktreeTargetDir: "worktrees",
			branchName:             "feature-branch",
			expected:               "/home/user/worktrees/feature-branch",
			description:            "homePath doesn't end with '/' and configWorktreeTargetDir doesn't end with '/'",
		},
		{
			name:                    "empty configWorktreeTargetDir with trailing slash on homePath",
			homePath:                "/home/user/",
			configWorktreeTargetDir: "",
			branchName:             "feature-branch",
			expected:               "/home/user/feature-branch",
			description:            "configWorktreeTargetDir is empty, homePath ends with '/'",
		},
		{
			name:                    "empty configWorktreeTargetDir without trailing slash on homePath",
			homePath:                "/home/user",
			configWorktreeTargetDir: "",
			branchName:             "feature-branch",
			expected:               "/home/user/feature-branch",
			description:            "configWorktreeTargetDir is empty, homePath doesn't end with '/'",
		},
		{
			name:                    "multiple slashes in paths",
			homePath:                "/home/user//",
			configWorktreeTargetDir: "//worktrees//",
			branchName:             "feature-branch",
			expected:               "/home/user/worktrees/feature-branch",
			description:            "paths with multiple slashes (should be cleaned up by filepath.Join)",
		},
		{
			name:                    "branch name with special characters",
			homePath:                "/home/user",
			configWorktreeTargetDir: "worktrees",
			branchName:             "feature/special-branch_123",
			expected:               "/home/user/worktrees/feature/special-branch_123",
			description:            "branch name containing special characters and slashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildConfigWorktreeDir(tt.homePath, tt.configWorktreeTargetDir, tt.branchName)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
} 