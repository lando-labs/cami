package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBackup(t *testing.T) {
	t.Run("create backup of directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test directory with some files
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(targetDir, "file1.txt"), []byte("content1"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(targetDir, "file2.txt"), []byte("content2"), 0644))

		// Create subdirectory
		subDir := filepath.Join(targetDir, "subdir")
		require.NoError(t, os.Mkdir(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("content3"), 0644))

		// Create backup
		backupPath, err := CreateBackup(targetDir)

		require.NoError(t, err)
		assert.NotEmpty(t, backupPath)
		assert.Contains(t, backupPath, BackupPrefix)

		// Verify backup exists
		assert.DirExists(t, backupPath)

		// Verify all files were copied
		assert.FileExists(t, filepath.Join(backupPath, "file1.txt"))
		assert.FileExists(t, filepath.Join(backupPath, "file2.txt"))
		assert.DirExists(t, filepath.Join(backupPath, "subdir"))
		assert.FileExists(t, filepath.Join(backupPath, "subdir", "file3.txt"))

		// Verify file contents
		content, err := os.ReadFile(filepath.Join(backupPath, "file1.txt"))
		require.NoError(t, err)
		assert.Equal(t, "content1", string(content))
	})

	t.Run("error on non-existent target", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistent := filepath.Join(tmpDir, "nonexistent")

		_, err := CreateBackup(nonExistent)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("error when target is file not directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "file.txt")
		require.NoError(t, os.WriteFile(testFile, []byte("content"), 0644))

		_, err := CreateBackup(testFile)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})

	t.Run("backup path is in same parent directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		backupPath, err := CreateBackup(targetDir)

		require.NoError(t, err)
		assert.Equal(t, tmpDir, filepath.Dir(backupPath))
	})
}

func TestListBackups(t *testing.T) {
	t.Run("list multiple backups", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create multiple backups with different timestamps
		backup1 := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		backup2 := filepath.Join(tmpDir, BackupPrefix+"20250102-120000")
		backup3 := filepath.Join(tmpDir, BackupPrefix+"20250103-120000")

		require.NoError(t, os.Mkdir(backup1, 0755))
		require.NoError(t, os.Mkdir(backup2, 0755))
		require.NoError(t, os.Mkdir(backup3, 0755))

		backups, err := ListBackups(targetDir)

		require.NoError(t, err)
		assert.Len(t, backups, 3)

		// Should be sorted newest first
		assert.Contains(t, backups[0].Path, "20250103")
		assert.Contains(t, backups[1].Path, "20250102")
		assert.Contains(t, backups[2].Path, "20250101")
	})

	t.Run("return empty list when no backups exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		backups, err := ListBackups(targetDir)

		require.NoError(t, err)
		assert.Empty(t, backups)
	})

	t.Run("ignore non-backup directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create backup directory
		backup := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		require.NoError(t, os.Mkdir(backup, 0755))

		// Create non-backup directories
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "other-dir"), 0755))
		require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "another-dir"), 0755))

		backups, err := ListBackups(targetDir)

		require.NoError(t, err)
		assert.Len(t, backups, 1)
		assert.Contains(t, backups[0].Path, BackupPrefix)
	})

	t.Run("calculate backup sizes", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create backup with files
		backup := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		require.NoError(t, os.Mkdir(backup, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(backup, "file.txt"), []byte("12345"), 0644))

		backups, err := ListBackups(targetDir)

		require.NoError(t, err)
		assert.Len(t, backups, 1)
		assert.Equal(t, int64(5), backups[0].SizeBytes)
	})
}

func TestAnalyzeArchive(t *testing.T) {
	t.Run("analyze archive with backups", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create backups
		backup1 := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		backup2 := filepath.Join(tmpDir, BackupPrefix+"20250102-120000")

		require.NoError(t, os.Mkdir(backup1, 0755))
		require.NoError(t, os.Mkdir(backup2, 0755))

		// Add files to backups
		require.NoError(t, os.WriteFile(filepath.Join(backup1, "file.txt"), []byte("12345"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(backup2, "file.txt"), []byte("1234567890"), 0644))

		analysis, err := AnalyzeArchive(targetDir)

		require.NoError(t, err)
		assert.Equal(t, 2, analysis.TotalBackups)
		assert.Equal(t, int64(15), analysis.TotalSizeBytes) // 5 + 10
		assert.Len(t, analysis.Backups, 2)

		// Check timestamps
		assert.True(t, analysis.NewestBackup.After(analysis.OldestBackup) || analysis.NewestBackup.Equal(analysis.OldestBackup))
	})

	t.Run("analyze empty archive", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		analysis, err := AnalyzeArchive(targetDir)

		require.NoError(t, err)
		assert.Equal(t, 0, analysis.TotalBackups)
		assert.Equal(t, int64(0), analysis.TotalSizeBytes)
		assert.Empty(t, analysis.Backups)
	})
}

func TestCleanupBackups(t *testing.T) {
	t.Run("cleanup old backups keeping recent N", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create 5 backups
		for i := 1; i <= 5; i++ {
			backup := filepath.Join(tmpDir, BackupPrefix+time.Now().Add(time.Duration(i)*time.Hour).Format("20060102-150405"))
			require.NoError(t, os.Mkdir(backup, 0755))
			require.NoError(t, os.WriteFile(filepath.Join(backup, "file.txt"), []byte("content"), 0644))
			time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		}

		// Keep only 3 most recent
		result, err := CleanupBackups(targetDir, CleanupOptions{KeepRecent: 3})

		require.NoError(t, err)
		assert.Equal(t, 2, result.RemovedCount)
		assert.Greater(t, result.FreedBytes, int64(0))
		assert.Len(t, result.KeptBackups, 3)

		// Verify only 3 backups remain
		backups, err := ListBackups(targetDir)
		require.NoError(t, err)
		assert.Len(t, backups, 3)
	})

	t.Run("no cleanup when backups below threshold", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create 2 backups
		backup1 := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		backup2 := filepath.Join(tmpDir, BackupPrefix+"20250102-120000")
		require.NoError(t, os.Mkdir(backup1, 0755))
		require.NoError(t, os.Mkdir(backup2, 0755))

		// Try to keep 3 (but only 2 exist)
		result, err := CleanupBackups(targetDir, CleanupOptions{KeepRecent: 3})

		require.NoError(t, err)
		assert.Equal(t, 0, result.RemovedCount)
		assert.Equal(t, int64(0), result.FreedBytes)
		assert.Len(t, result.KeptBackups, 2)

		// Verify both backups still exist
		backups, err := ListBackups(targetDir)
		require.NoError(t, err)
		assert.Len(t, backups, 2)
	})

	t.Run("use default keep count when not specified", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create 5 backups
		for i := 1; i <= 5; i++ {
			backup := filepath.Join(tmpDir, BackupPrefix+time.Now().Add(time.Duration(i)*time.Hour).Format("20060102-150405"))
			require.NoError(t, os.Mkdir(backup, 0755))
			time.Sleep(1 * time.Millisecond)
		}

		// Use default (should be 3)
		result, err := CleanupBackups(targetDir, CleanupOptions{KeepRecent: 0})

		require.NoError(t, err)
		assert.Equal(t, 2, result.RemovedCount)
		assert.Len(t, result.KeptBackups, DefaultKeepRecent)
	})
}

func TestRestoreFromBackup(t *testing.T) {
	t.Run("restore backup to target location", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create backup directory with files
		backupPath := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		require.NoError(t, os.Mkdir(backupPath, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(backupPath, "file1.txt"), []byte("backup content"), 0644))

		subDir := filepath.Join(backupPath, "subdir")
		require.NoError(t, os.Mkdir(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("nested content"), 0644))

		// Target location
		targetPath := filepath.Join(tmpDir, "restored")

		// Restore
		err := RestoreFromBackup(backupPath, targetPath)

		require.NoError(t, err)
		assert.DirExists(t, targetPath)

		// Verify files were restored
		assert.FileExists(t, filepath.Join(targetPath, "file1.txt"))
		assert.DirExists(t, filepath.Join(targetPath, "subdir"))
		assert.FileExists(t, filepath.Join(targetPath, "subdir", "file2.txt"))

		// Verify content
		content, err := os.ReadFile(filepath.Join(targetPath, "file1.txt"))
		require.NoError(t, err)
		assert.Equal(t, "backup content", string(content))
	})

	t.Run("restore replaces existing target", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create backup
		backupPath := filepath.Join(tmpDir, BackupPrefix+"20250101-120000")
		require.NoError(t, os.Mkdir(backupPath, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(backupPath, "new.txt"), []byte("new content"), 0644))

		// Create existing target with different content
		targetPath := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetPath, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(targetPath, "old.txt"), []byte("old content"), 0644))

		// Restore (should replace)
		err := RestoreFromBackup(backupPath, targetPath)

		require.NoError(t, err)

		// Verify old file is gone, new file exists
		assert.NoFileExists(t, filepath.Join(targetPath, "old.txt"))
		assert.FileExists(t, filepath.Join(targetPath, "new.txt"))
	})

	t.Run("error on non-existent backup", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupPath := filepath.Join(tmpDir, BackupPrefix+"nonexistent")
		targetPath := filepath.Join(tmpDir, "target")

		err := RestoreFromBackup(backupPath, targetPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("error on invalid backup name", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create directory without backup prefix
		invalidBackup := filepath.Join(tmpDir, "not-a-backup")
		require.NoError(t, os.Mkdir(invalidBackup, 0755))

		targetPath := filepath.Join(tmpDir, "target")

		err := RestoreFromBackup(invalidBackup, targetPath)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid backup")
	})
}

func TestShouldSuggestCleanup(t *testing.T) {
	t.Run("suggest cleanup when at threshold", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create exactly CleanupThreshold backups
		for i := 1; i <= CleanupThreshold; i++ {
			backup := filepath.Join(tmpDir, BackupPrefix+time.Now().Add(time.Duration(i)*time.Minute).Format("20060102-150405"))
			require.NoError(t, os.Mkdir(backup, 0755))
			time.Sleep(1 * time.Millisecond)
		}

		shouldCleanup, err := ShouldSuggestCleanup(targetDir)

		require.NoError(t, err)
		assert.True(t, shouldCleanup)
	})

	t.Run("do not suggest when below threshold", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := filepath.Join(tmpDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0755))

		// Create fewer than threshold
		for i := 1; i < CleanupThreshold; i++ {
			backup := filepath.Join(tmpDir, BackupPrefix+time.Now().Add(time.Duration(i)*time.Minute).Format("20060102-150405"))
			require.NoError(t, os.Mkdir(backup, 0755))
			time.Sleep(1 * time.Millisecond)
		}

		shouldCleanup, err := ShouldSuggestCleanup(targetDir)

		require.NoError(t, err)
		assert.False(t, shouldCleanup)
	})
}

func TestCopyDir(t *testing.T) {
	t.Run("copy directory with nested structure", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source directory with nested structure
		srcDir := filepath.Join(tmpDir, "src")
		require.NoError(t, os.Mkdir(srcDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644))

		subDir := filepath.Join(srcDir, "subdir")
		require.NoError(t, os.Mkdir(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("content2"), 0644))

		// Copy to destination
		dstDir := filepath.Join(tmpDir, "dst")
		err := copyDir(srcDir, dstDir)

		require.NoError(t, err)
		assert.DirExists(t, dstDir)
		assert.FileExists(t, filepath.Join(dstDir, "file1.txt"))
		assert.DirExists(t, filepath.Join(dstDir, "subdir"))
		assert.FileExists(t, filepath.Join(dstDir, "subdir", "file2.txt"))

		// Verify contents
		content, err := os.ReadFile(filepath.Join(dstDir, "file1.txt"))
		require.NoError(t, err)
		assert.Equal(t, "content1", string(content))
	})
}

func TestCalculateDirSize(t *testing.T) {
	t.Run("calculate directory size", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create directory with files
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("12345"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("1234567890"), 0644))

		size, err := calculateDirSize(tmpDir)

		require.NoError(t, err)
		assert.Equal(t, int64(15), size) // 5 + 10
	})

	t.Run("calculate size with subdirectories", func(t *testing.T) {
		tmpDir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("12345"), 0644))

		subDir := filepath.Join(tmpDir, "subdir")
		require.NoError(t, os.Mkdir(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("1234567890"), 0644))

		size, err := calculateDirSize(tmpDir)

		require.NoError(t, err)
		assert.Equal(t, int64(15), size)
	})
}
