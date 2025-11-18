package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// BackupInfo represents a backup directory
type BackupInfo struct {
	Path      string
	Timestamp time.Time
	SizeBytes int64
}

// ArchiveAnalysis represents backup state
type ArchiveAnalysis struct {
	TotalBackups   int
	TotalSizeBytes int64
	OldestBackup   time.Time
	NewestBackup   time.Time
	Backups        []BackupInfo
}

// CleanupOptions specifies what to keep
type CleanupOptions struct {
	KeepRecent int // Number of recent backups to keep
}

// CleanupResult represents the outcome
type CleanupResult struct {
	RemovedCount int
	FreedBytes   int64
	KeptBackups  []string
}

const (
	// BackupPrefix is the prefix for backup directories
	BackupPrefix = ".cami-backup-"

	// DefaultKeepRecent is the default number of backups to keep
	DefaultKeepRecent = 3

	// CleanupThreshold is the number of backups that triggers cleanup suggestion
	CleanupThreshold = 10
)

// CreateBackup creates a backup of the target directory
func CreateBackup(targetPath string) (string, error) {
	// Validate target exists
	info, err := os.Stat(targetPath)
	if err != nil {
		return "", fmt.Errorf("target path does not exist: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("target path is not a directory: %s", targetPath)
	}

	// Generate backup directory name with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("%s%s", BackupPrefix, timestamp)

	// Backup goes in the same parent directory as target
	parentDir := filepath.Dir(targetPath)
	backupPath := filepath.Join(parentDir, backupName)

	// Copy entire directory
	if err := copyDir(targetPath, backupPath); err != nil {
		// Clean up partial backup on failure
		os.RemoveAll(backupPath)
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return backupPath, nil
}

// ListBackups lists all backup directories for a given target path
func ListBackups(targetPath string) ([]BackupInfo, error) {
	parentDir := filepath.Dir(targetPath)

	// Read parent directory
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		// Check if directory starts with backup prefix
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), BackupPrefix) {
			continue
		}

		backupPath := filepath.Join(parentDir, entry.Name())

		// Get directory info
		info, err := os.Stat(backupPath)
		if err != nil {
			continue
		}

		// Parse timestamp from backup name
		timestampStr := strings.TrimPrefix(entry.Name(), BackupPrefix)
		timestamp, err := time.Parse("20060102-150405", timestampStr)
		if err != nil {
			// If we can't parse timestamp, use modification time
			timestamp = info.ModTime()
		}

		// Calculate directory size
		size, _ := calculateDirSize(backupPath)

		backups = append(backups, BackupInfo{
			Path:      backupPath,
			Timestamp: timestamp,
			SizeBytes: size,
		})
	}

	// Sort by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	return backups, nil
}

// AnalyzeArchive analyzes the backup state for a target path
func AnalyzeArchive(targetPath string) (*ArchiveAnalysis, error) {
	backups, err := ListBackups(targetPath)
	if err != nil {
		return nil, err
	}

	if len(backups) == 0 {
		return &ArchiveAnalysis{
			TotalBackups:   0,
			TotalSizeBytes: 0,
			Backups:        []BackupInfo{},
		}, nil
	}

	analysis := &ArchiveAnalysis{
		TotalBackups: len(backups),
		Backups:      backups,
	}

	// Calculate total size
	for _, backup := range backups {
		analysis.TotalSizeBytes += backup.SizeBytes
	}

	// Find oldest and newest
	analysis.NewestBackup = backups[0].Timestamp
	analysis.OldestBackup = backups[len(backups)-1].Timestamp

	return analysis, nil
}

// CleanupBackups removes old backups, keeping the most recent N
func CleanupBackups(targetPath string, options CleanupOptions) (*CleanupResult, error) {
	backups, err := ListBackups(targetPath)
	if err != nil {
		return nil, err
	}

	keepCount := options.KeepRecent
	if keepCount <= 0 {
		keepCount = DefaultKeepRecent
	}

	// If we have fewer backups than we want to keep, nothing to do
	if len(backups) <= keepCount {
		keptPaths := make([]string, len(backups))
		for i, b := range backups {
			keptPaths[i] = b.Path
		}
		return &CleanupResult{
			RemovedCount: 0,
			FreedBytes:   0,
			KeptBackups:  keptPaths,
		}, nil
	}

	result := &CleanupResult{}

	// Keep the first N (newest) backups
	for i := 0; i < keepCount; i++ {
		result.KeptBackups = append(result.KeptBackups, backups[i].Path)
	}

	// Remove the rest
	for i := keepCount; i < len(backups); i++ {
		backup := backups[i]
		if err := os.RemoveAll(backup.Path); err != nil {
			return nil, fmt.Errorf("failed to remove backup %s: %w", backup.Path, err)
		}
		result.RemovedCount++
		result.FreedBytes += backup.SizeBytes
	}

	return result, nil
}

// RestoreFromBackup restores a backup to the target location
func RestoreFromBackup(backupPath string, targetPath string) error {
	// Validate backup exists
	info, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("backup does not exist: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("backup is not a directory: %s", backupPath)
	}

	// Validate it's actually a backup directory
	backupName := filepath.Base(backupPath)
	if !strings.HasPrefix(backupName, BackupPrefix) {
		return fmt.Errorf("not a valid backup directory: %s", backupPath)
	}

	// Remove target if it exists
	if _, err := os.Stat(targetPath); err == nil {
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("failed to remove existing target: %w", err)
		}
	}

	// Copy backup to target location
	if err := copyDir(backupPath, targetPath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	return nil
}

// ShouldSuggestCleanup returns true if cleanup should be suggested to user
func ShouldSuggestCleanup(targetPath string) (bool, error) {
	backups, err := ListBackups(targetPath)
	if err != nil {
		return false, err
	}

	return len(backups) >= CleanupThreshold, nil
}

// copyDir recursively copies a directory
func copyDir(src string, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// calculateDirSize calculates the total size of a directory
func calculateDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
