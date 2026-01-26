package git

import (
	"testing"

	com "github.com/garrettkrohn/treekanga/common"
	"github.com/stretchr/testify/assert"
)

func TestDetermineBranchArguments(t *testing.T) {
	// Create a RealGitAdapter instance for testing
	git := &RealGitAdapter{}

	t.Run("Case 1a: New branch exists locally", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  true,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: true,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"feature-branch"}

		assert.Equal(t, expected, result, "Should return just the branch name when new branch exists locally")
	})

	t.Run("Case 1b: New branch exists remotely", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: true,
				BaseBranchExistsLocally: true,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"feature-branch"}

		assert.Equal(t, expected, result, "Should return just the branch name when new branch exists remotely")
	})

	t.Run("Case 1c: New branch exists both locally and remotely", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  true,
				NewBranchExistsRemotely: true,
				BaseBranchExistsLocally: true,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"feature-branch"}

		assert.Equal(t, expected, result, "Should return just the branch name when new branch exists both locally and remotely")
	})

	t.Run("Case 2a: Base branch exists locally and should pull", func(t *testing.T) {
		pullFlag := true
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: true,
			},
			Flags: com.AddCmdFlags{
				Pull: &pullFlag,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "feature-branch", "origin/main", "--no-track"}

		assert.Equal(t, expected, result, "Should create new branch from remote version when pull flag is set")
	})

	t.Run("Case 2b: Base branch exists locally and should not pull", func(t *testing.T) {
		pullFlag := false
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: true,
			},
			Flags: com.AddCmdFlags{
				Pull: &pullFlag,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "feature-branch", "main"}

		assert.Equal(t, expected, result, "Should create new branch from local version when pull flag is false")
	})

	t.Run("Case 2c: Base branch exists locally with nil pull flag", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: true,
			},
			Flags: com.AddCmdFlags{
				Pull: nil, // ShouldPull() will return false
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "feature-branch", "main"}

		assert.Equal(t, expected, result, "Should create new branch from local version when pull flag is nil")
	})

	t.Run("Case 3: Base branch only exists remotely", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature-branch",
				BaseBranchName:          "develop",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: false,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "feature-branch", "origin/develop", "--no-track"}

		assert.Equal(t, expected, result, "Should create new branch from remote when base branch only exists remotely")
	})

	t.Run("Edge case: Empty branch names", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "",
				BaseBranchName:          "",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: false,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "", "origin/", "--no-track"}

		assert.Equal(t, expected, result, "Should handle empty branch names gracefully")
	})
}

// TestDetermineBranchArgumentsIntegration tests the function with more realistic scenarios
func TestDetermineBranchArgumentsIntegration(t *testing.T) {
	git := &RealGitAdapter{}

	t.Run("Typical new feature branch from main", func(t *testing.T) {
		pullFlag := false
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature/user-authentication",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: true,
			},
			Flags: com.AddCmdFlags{
				Pull: &pullFlag,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "feature/user-authentication", "main"}

		assert.Equal(t, expected, result, "Should create feature branch from local main")
	})

	t.Run("Hotfix branch with pull from remote", func(t *testing.T) {
		pullFlag := true
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "hotfix/critical-bug",
				BaseBranchName:          "production",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: false,
				BaseBranchExistsLocally: true,
			},
			Flags: com.AddCmdFlags{
				Pull: &pullFlag,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"-b", "hotfix/critical-bug", "origin/production", "--no-track"}

		assert.Equal(t, expected, result, "Should create hotfix branch from remote production")
	})

	t.Run("Checkout existing remote branch", func(t *testing.T) {
		config := &com.AddConfig{
			GitInfo: com.GitInfo{
				NewBranchName:           "feature/existing-feature",
				BaseBranchName:          "main",
				NewBranchExistsLocally:  false,
				NewBranchExistsRemotely: true,
				BaseBranchExistsLocally: true,
			},
		}

		result := git.determineBranchArguments(config)
		expected := []string{"feature/existing-feature"}

		assert.Equal(t, expected, result, "Should checkout existing remote branch")
	})
}
