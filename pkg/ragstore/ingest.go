package ragstore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var textExtensions = map[string]struct{}{
	".txt": {}, ".md": {}, ".markdown": {}, ".rst": {}, ".json": {}, ".yaml": {}, ".yml": {}, ".csv": {},
	".html": {}, ".htm": {},
}

var htmlTagPattern = regexp.MustCompile(`(?s)<[^>]+>`)

// ReadOptions configures directory ingestion.
type ReadOptions struct {
	Recursive bool
	Glob      string
	Exclude   []string
}

// ReadSources loads textual content from files and directories.
func ReadSources(paths []string, opts ReadOptions) ([]SourceDoc, error) {
	var docs []SourceDoc
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		st, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if st.IsDir() {
			found, err := readDirFiltered(p, opts)
			if err != nil {
				return nil, err
			}
			docs = append(docs, found...)
			continue
		}
		if !matchReadFilters(p, opts) {
			continue
		}
		doc, err := readFile(p)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// ReadStdin loads one document from stdin.
func ReadStdin(r io.Reader) (SourceDoc, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return SourceDoc{}, err
	}
	return SourceDoc{Path: "stdin", Content: string(data)}, nil
}

type SourceDoc struct {
	Path    string
	Content string
}

func readDirFiltered(dir string, opts ReadOptions) ([]SourceDoc, error) {
	var docs []SourceDoc
	walkFn := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != dir && !opts.Recursive {
				return filepath.SkipDir
			}
			if path != dir && matchExcludeDir(path, opts.Exclude) {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTextFile(path) || !matchReadFilters(path, opts) {
			return nil
		}
		doc, err := readFile(path)
		if err != nil {
			return err
		}
		docs = append(docs, doc)
		return nil
	}
	if err := filepath.WalkDir(dir, walkFn); err != nil {
		return nil, err
	}
	return docs, nil
}

func readFile(path string) (SourceDoc, error) {
	if !isTextFile(path) {
		return SourceDoc{}, fmt.Errorf("unsupported file type %q", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return SourceDoc{}, err
	}
	content := string(data)
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".html" || ext == ".htm" {
		content = stripHTML(content)
	}
	return SourceDoc{Path: path, Content: content}, nil
}

func stripHTML(s string) string {
	s = htmlTagPattern.ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := textExtensions[ext]
	return ok
}

func matchReadFilters(path string, opts ReadOptions) bool {
	if len(opts.Exclude) > 0 {
		for _, ex := range opts.Exclude {
			ex = strings.TrimSpace(ex)
			if ex == "" {
				continue
			}
			if strings.Contains(path, ex) {
				return false
			}
			if matched, _ := filepath.Match(ex, filepath.Base(path)); matched {
				return false
			}
		}
	}
	if opts.Glob == "" {
		return true
	}
	base := filepath.Base(path)
	matched, err := filepath.Match(opts.Glob, base)
	return err == nil && matched
}

func matchExcludeDir(path string, exclude []string) bool {
	for _, ex := range exclude {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		if strings.Contains(path, ex) {
			return true
		}
	}
	return false
}
