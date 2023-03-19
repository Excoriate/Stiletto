package git

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/logger"
	"os"
	"path/filepath"
)

// FindGitRootDir finds the root directory of a git repository
func FindGitRootDir(l *logger.StilettoLog, path string, limit int) (string, error) {
	levelsChecked := 0

	for {
		if levelsChecked > limit {
			break
		}

		gitPath := filepath.Join(path, ".git")
		l.LogDebug("Checking if %s is a git repository", gitPath)
		info, err := os.Stat(gitPath)

		if err == nil && info.IsDir() {
			return path, nil
		}

		newPath := filepath.Dir(path)
		if newPath == path {
			break
		}
		path = newPath
		levelsChecked++
	}

	l.LogDebug(fmt.Sprintf("Could not find git repository in %d levels", limit))
	return "", fmt.Errorf("could not find git repository in %d levels", limit)
}

// IsAGitRepository checks if a path is a git repository
func IsAGitRepository(l *logger.StilettoLog, path string) bool {
	_, err := FindGitRootDir(l, path, 0)
	return err == nil
}
