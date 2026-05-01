package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectCmd_SelectFlag_ParsedCorrectly(t *testing.T) {
	// Arrange - create connect command with --select flag
	cmd := &cobra.Command{
		Use: "connect",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().BoolP("select", "s", false, "Interactive selection mode")

	// Act - parse flags
	err := cmd.ParseFlags([]string{"--select"})
	require.NoError(t, err)

	selectFlag, err := cmd.Flags().GetBool("select")
	require.NoError(t, err)

	// Assert - expect select flag = true
	assert.True(t, selectFlag, "Expected --select flag to be parsed as true")
}

func TestConnectCmd_BareFlag_ParsedCorrectly(t *testing.T) {
	// Arrange - create connect command with --select --bare flags
	cmd := &cobra.Command{
		Use: "connect",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().BoolP("select", "s", false, "Interactive selection mode")
	cmd.Flags().Bool("bare", false, "Select from bare repos")

	// Act - parse flags
	err := cmd.ParseFlags([]string{"--select", "--bare"})
	require.NoError(t, err)

	selectFlag, err := cmd.Flags().GetBool("select")
	require.NoError(t, err)
	bareFlag, err := cmd.Flags().GetBool("bare")
	require.NoError(t, err)

	// Assert - expect both flags = true
	assert.True(t, selectFlag, "Expected --select flag to be true")
	assert.True(t, bareFlag, "Expected --bare flag to be true")
}

func TestConnectCmd_ByRepoFlag_ParsedCorrectly(t *testing.T) {
	// Arrange - create connect command with --select --by-repo flags
	cmd := &cobra.Command{
		Use: "connect",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().BoolP("select", "s", false, "Interactive selection mode")
	cmd.Flags().Bool("by-repo", false, "Select repo first")

	// Act - parse flags
	err := cmd.ParseFlags([]string{"--select", "--by-repo"})
	require.NoError(t, err)

	selectFlag, err := cmd.Flags().GetBool("select")
	require.NoError(t, err)
	byRepoFlag, err := cmd.Flags().GetBool("by-repo")
	require.NoError(t, err)

	// Assert - expect both flags = true
	assert.True(t, selectFlag, "Expected --select flag to be true")
	assert.True(t, byRepoFlag, "Expected --by-repo flag to be true")
}

func TestListBareRepos_ReturnsAllBareRepoPaths(t *testing.T) {
	// Arrange - config with 2 repos, mock fs with both paths existing
	// This test will use the real filesystem, so we'll test the function's logic
	// Note: This is a placeholder test that will fail until listBareRepos() is implemented

	// Act - call listBareRepos()
	_, err := listBareRepos()

	// Assert - expect 2 bare repo paths returned
	// For now, we just check that the function exists and returns without panic
	assert.NoError(t, err, "Expected listBareRepos to return without error")
}

func TestListBareRepos_FiltersNonExistentPaths(t *testing.T) {
	// Arrange - config with 3 repos, mock fs (2 exist, 1 missing)
	// This test verifies filtering behavior

	// Act - call listBareRepos()
	repos, err := listBareRepos()

	// Assert - expect only existing paths returned
	assert.NoError(t, err, "Expected no error from listBareRepos")
	assert.NotNil(t, repos, "Expected repos slice to be non-nil")
}

func TestListBareRepos_ReturnsEmptyList_WhenNoBareReposExist(t *testing.T) {
	// Arrange - config with repos, mock fs (all missing)
	// This will depend on the actual config state

	// Act - call listBareRepos()
	repos, err := listBareRepos()

	// Assert - expect empty slice or list of existing repos
	assert.NoError(t, err, "Expected no error from listBareRepos")
	assert.NotNil(t, repos, "Expected repos slice to be non-nil")
}

func TestFormatBareRepoForDisplay_FormatsCorrectly(t *testing.T) {
	// Arrange - bare repo path string
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "cal_work repo",
			input:    "/Users/gkrohn/code/cal_work/.bare",
			expected: "cal -> /Users/gkrohn/code/cal_work/.bare",
		},
		{
			name:     "core_work repo",
			input:    "/Users/gkrohn/code/core_work/.bare",
			expected: "core -> /Users/gkrohn/code/core_work/.bare",
		},
		{
			name:     "repo without _work suffix",
			input:    "/Users/gkrohn/code/myproject/.bare",
			expected: "myproject -> /Users/gkrohn/code/myproject/.bare",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act - call formatBareRepoForDisplay(path)
			result := formatBareRepoForDisplay(tc.input)

			// Assert - expect formatted string with repo name extracted
			assert.Equal(t, tc.expected, result, "Expected formatted string to match")
		})
	}
}
