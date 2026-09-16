package remotemedia

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// uploadEntry is one file to upload, named relative to the staging session.
type uploadEntry struct {
	LocalPath string
	BlobName  string
	RootPath  string // set when LocalPath belongs to a staged directory
}

// dirUpload tracks a local directory mapped to a server staging subdirectory.
type dirUpload struct {
	LocalDir string
	DirName  string
}

// BuildUploadPlan expands files and directories into a flat upload list.
func BuildUploadPlan(paths []string) ([]uploadEntry, []dirUpload, error) {
	seen := make(map[string]struct{})
	var files []uploadEntry
	var dirs []dirUpload

	for _, raw := range paths {
		p := strings.TrimSpace(raw)
		if p == "" || isRemoteURL(p) {
			continue
		}
		p = filepath.Clean(p)
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}

		st, err := os.Stat(p)
		if err != nil {
			continue
		}
		if st.IsDir() {
			dirName := filepath.Base(p)
			if dirName == "." || dirName == string(os.PathSeparator) {
				dirName = "data"
			}
			dirs = append(dirs, dirUpload{LocalDir: p, DirName: dirName})
			err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				rel, err := filepath.Rel(p, path)
				if err != nil {
					return err
				}
				files = append(files, uploadEntry{
					LocalPath: path,
					BlobName:  filepath.Join(dirName, rel),
					RootPath:  p,
				})
				return nil
			})
			if err != nil {
				return nil, nil, fmt.Errorf("walk %s: %w", p, err)
			}
			continue
		}
		files = append(files, uploadEntry{
			LocalPath: p,
			BlobName:  fmt.Sprintf("file-%d_%s", len(files), filepath.Base(p)),
		})
	}
	return files, dirs, nil
}

// ServerDirForLocal returns the server-side directory for a staged local folder.
func ServerDirForLocal(stagingRoot, localDir string) string {
	dirName := filepath.Base(filepath.Clean(localDir))
	if dirName == "." || dirName == string(os.PathSeparator) {
		dirName = "data"
	}
	return filepath.Join(stagingRoot, dirName)
}
