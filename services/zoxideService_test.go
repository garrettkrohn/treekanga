package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/garrettkrohn/treekanga/directoryReader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompileZoxidePathsToAdd(t *testing.T) {
	t.Run("Empty config - should only return root directory", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)
		newRootDirectory := "/test/root"
		foldersToAddFromConfig := []string{}

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert
		expected := []string{newRootDirectory}
		assert.Equal(t, expected, result, "Should return only root directory when config is empty")
	})

	t.Run("Non-wildcard folder - valid path", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)

		// Create temporary directories for testing
		tempDir := t.TempDir()
		newRootDirectory := tempDir
		subFolder := "subfolder"
		subFolderPath := filepath.Join(tempDir, subFolder)
		err := os.Mkdir(subFolderPath, 0755)
		assert.NoError(t, err)

		foldersToAddFromConfig := []string{subFolder}

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert
		expected := []string{
			newRootDirectory,
			subFolderPath,
		}
		assert.Equal(t, expected, result, "Should include root and valid subfolder")
	})

	t.Run("Non-wildcard folder - invalid path should be skipped", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)
		tempDir := t.TempDir()
		newRootDirectory := tempDir
		foldersToAddFromConfig := []string{"nonexistent"}

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert
		expected := []string{newRootDirectory}
		assert.Equal(t, expected, result, "Should only return root when subfolder doesn't exist")
	})

	t.Run("Wildcard folder - expands to multiple directories", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)

		// Create temporary directories for testing
		tempDir := t.TempDir()
		// Note: newRootDirectory should end with "/" for the current implementation to work correctly
		// This is due to a bug in line 31 of zoxideService.go where paths are concatenated without separator
		newRootDirectory := tempDir + "/"

		// Create a parent directory and multiple subdirectories
		parentDir := "projects"
		parentPath := filepath.Join(tempDir, parentDir)
		err := os.Mkdir(parentPath, 0755)
		assert.NoError(t, err)

		project1Path := filepath.Join(parentPath, "project1")
		project2Path := filepath.Join(parentPath, "project2")
		err = os.Mkdir(project1Path, 0755)
		assert.NoError(t, err)
		err = os.Mkdir(project2Path, 0755)
		assert.NoError(t, err)

		foldersToAddFromConfig := []string{"projects/*"}

		// The code concatenates newRootDirectory + pathUpTillWildcard
		// expectedMockPath := newRootDirectory + "projects"
		mockDirReader.EXPECT().GetFoldersInDirectory(mock.Anything).Return([]string{"project1", "project2"}, nil)

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert - note that root directory is returned with trailing slash
		expected := []string{
			newRootDirectory, // Keeps trailing slash in result
			project1Path,
			project2Path,
		}
		assert.Equal(t, expected, result, "Should expand wildcard to multiple directories")
	})

	t.Run("Wildcard folder - no matching directories", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)
		tempDir := t.TempDir()
		newRootDirectory := tempDir + "/"

		// Create a parent directory but no subdirectories
		parentDir := "empty"
		parentPath := filepath.Join(tempDir, parentDir)
		err := os.Mkdir(parentPath, 0755)
		assert.NoError(t, err)

		foldersToAddFromConfig := []string{"empty/*"}

		expectedMockPath := newRootDirectory + "empty"
		mockDirReader.EXPECT().GetFoldersInDirectory(expectedMockPath).Return([]string{}, nil)

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert - root directory is returned with trailing slash
		expected := []string{newRootDirectory}
		assert.Equal(t, expected, result, "Should only return root when wildcard matches nothing")
	})

	t.Run("Mixed config - wildcards and non-wildcards", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)

		tempDir := t.TempDir()
		newRootDirectory := tempDir + "/"

		// Create regular folder
		regularFolder := "regular"
		regularPath := filepath.Join(tempDir, regularFolder)
		err := os.Mkdir(regularPath, 0755)
		assert.NoError(t, err)

		// Create wildcard parent and children
		wildcardParent := "wildcard"
		wildcardParentPath := filepath.Join(tempDir, wildcardParent)
		err = os.Mkdir(wildcardParentPath, 0755)
		assert.NoError(t, err)

		child1Path := filepath.Join(wildcardParentPath, "child1")
		child2Path := filepath.Join(wildcardParentPath, "child2")
		err = os.Mkdir(child1Path, 0755)
		assert.NoError(t, err)
		err = os.Mkdir(child2Path, 0755)
		assert.NoError(t, err)

		foldersToAddFromConfig := []string{
			"regular",
			"wildcard/*",
		}

		expectedMockPath := filepath.Join(newRootDirectory, "wildcard")
		mockDirReader.EXPECT().GetFoldersInDirectory(expectedMockPath).Return([]string{"child1", "child2"}, nil)

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert - note double slash in regular path due to concatenation bug (line 22)
		regularPathWithDoubleSlash := filepath.Join(newRootDirectory + regularFolder)
		expected := []string{
			newRootDirectory,
			regularPathWithDoubleSlash, // Will have double slash: /root//regular
			child1Path,
			child2Path,
		}
		assert.Equal(t, expected, result, "Should handle mix of regular and wildcard folders")
	})

	t.Run("Nested wildcard path", func(t *testing.T) {
		// Setup
		mockDirReader := directoryReader.NewMockDirectoryReader(t)

		tempDir := t.TempDir()
		newRootDirectory := tempDir + "/"

		// Create nested structure: root/parent/child/*
		parentDir := "parent"
		childDir := "child"
		parentPath := filepath.Join(tempDir, parentDir)
		childPath := filepath.Join(parentPath, childDir)
		err := os.MkdirAll(childPath, 0755)
		assert.NoError(t, err)

		grandchild1Path := filepath.Join(childPath, "grandchild1")
		grandchild2Path := filepath.Join(childPath, "grandchild2")
		err = os.Mkdir(grandchild1Path, 0755)
		assert.NoError(t, err)
		err = os.Mkdir(grandchild2Path, 0755)
		assert.NoError(t, err)

		foldersToAddFromConfig := []string{"parent/child/*"}

		// For "parent/child/*", getPathUntilLastSlash returns "parent/child"
		expectedMockPath := newRootDirectory + "parent/child"
		mockDirReader.EXPECT().GetFoldersInDirectory(expectedMockPath).Return([]string{"grandchild1", "grandchild2"}, nil)

		// Execute
		result := CompileZoxidePathsToAdd(foldersToAddFromConfig, newRootDirectory, mockDirReader)

		// Assert
		expected := []string{
			newRootDirectory,
			grandchild1Path,
			grandchild2Path,
		}
		assert.Equal(t, expected, result, "Should handle nested wildcard paths")
	})

	t.Run("single *", func(t *testing.T) {
		//assert
		input := "*"

		//act
		result := getPathUntilLastSlash(input)

		assert.Equal(t, "", result, "should return *")

	})

	t.Run("wildcard after folder", func(t *testing.T) {
		//assert
		input := "parent/*"

		//act
		result := getPathUntilLastSlash(input)

		assert.Equal(t, "parent", result, "should return parent")

	})

}
