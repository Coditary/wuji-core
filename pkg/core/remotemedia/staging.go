package remotemedia

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
)

const stagingDir = ".wuji/staging"

// StagingSessionRoot returns the absolute session directory on the server.
func StagingSessionRoot(root, sessionID string) string {
	return filepath.Join(root, stagingDir, sanitizeSessionID(sessionID))
}

// WriteStaging saves uploaded blobs under <root>/.wuji/staging/<sessionID>/.
func WriteStaging(root, sessionID string, files []*wujiv1.FileBlob) ([]string, string, error) {
	if sessionID == "" {
		return nil, "", fmt.Errorf("session_id is required")
	}
	base := StagingSessionRoot(root, sessionID)
	if len(files) == 0 {
		return nil, base, nil
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return nil, "", err
	}
	out := make([]string, 0, len(files))
	for i, f := range files {
		name := safeBlobName(f.GetName(), i)
		dest := filepath.Join(base, name)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, "", err
		}
		raw, err := MaybeDecompress(f.GetData(), f.GetZstd())
		if err != nil {
			return nil, "", fmt.Errorf("decompress %s: %w", name, err)
		}
		if err := os.WriteFile(dest, raw, 0o644); err != nil {
			return nil, "", err
		}
		out = append(out, dest)
	}
	return out, base, nil
}

// ReadFiles reads files or directory trees for download to remote clients.
func ReadFiles(paths []string) ([]*wujiv1.FileBlob, error) {
	var out []*wujiv1.FileBlob
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		st, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", p, err)
		}
		if st.IsDir() {
			blobs, err := readDirTree(p)
			if err != nil {
				return nil, err
			}
			out = append(out, blobs...)
			continue
		}
		blob, err := readFileBlob(p, filepath.Base(p))
		if err != nil {
			return nil, err
		}
		out = append(out, blob)
	}
	return out, nil
}

func readDirTree(root string) ([]*wujiv1.FileBlob, error) {
	var out []*wujiv1.FileBlob
	base := filepath.Base(root)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		name := filepath.Join(base, rel)
		blob, err := readFileBlob(path, name)
		if err != nil {
			return err
		}
		out = append(out, blob)
		return nil
	})
	return out, err
}

func readFileBlob(path, name string) (*wujiv1.FileBlob, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	payload, zstdFlag, err := BlobData(data)
	if err != nil {
		return nil, err
	}
	return &wujiv1.FileBlob{Name: name, Data: payload, Zstd: zstdFlag}, nil
}

func safeBlobName(name string, index int) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Sprintf("file-%d", index)
	}
	name = filepath.Clean(strings.ReplaceAll(name, "..", ""))
	name = strings.TrimPrefix(name, string(os.PathSeparator))
	if name == "." || name == "" {
		return fmt.Sprintf("file-%d", index)
	}
	return name
}

// ResolvePaths turns relative paths into absolute paths under the project root.
func ResolvePaths(root string, paths []string) []string {
	if root == "" {
		root = "."
	}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !filepath.IsAbs(p) {
			p = filepath.Join(root, p)
		}
		out = append(out, filepath.Clean(p))
	}
	return out
}

func sanitizeSessionID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return "default"
	}
	var b strings.Builder
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}
