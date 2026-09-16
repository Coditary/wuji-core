package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LoadIndex reads the slim catalog index from disk.
func LoadIndex(root string) (*Index, error) {
	data, err := os.ReadFile(IndexPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("catalog not synced — run: wuji catalog sync")
		}
		return nil, fmt.Errorf("read catalog index: %w", err)
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse catalog index: %w", err)
	}
	if idx.Providers == nil {
		idx.Providers = map[string]ProviderEntry{}
	}
	if idx.Models == nil {
		idx.Models = map[string]ModelEntry{}
	}
	if idx.Labs == nil {
		idx.Labs = map[string]LabEntry{}
	}
	return &idx, nil
}

// LoadMeta reads sync metadata from disk.
func LoadMeta(root string) (*Meta, error) {
	data, err := os.ReadFile(MetaPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read catalog meta: %w", err)
	}
	var meta Meta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("parse catalog meta: %w", err)
	}
	return &meta, nil
}

// Save writes index.json and meta.json atomically under .wuji/catalog/.
func Save(root string, idx *Index, meta *Meta) error {
	if idx == nil || meta == nil {
		return fmt.Errorf("catalog index and meta are required")
	}
	dir := Dir(root)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create catalog dir: %w", err)
	}

	indexData, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("encode catalog index: %w", err)
	}
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("encode catalog meta: %w", err)
	}

	if err := writeAtomic(filepath.Join(dir, "index.json"), indexData); err != nil {
		return err
	}
	if err := writeAtomic(filepath.Join(dir, "meta.json"), metaData); err != nil {
		return err
	}
	return nil
}

func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".catalog-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("chmod %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close %s: %w", path, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		cleanup()
		return fmt.Errorf("rename %s: %w", path, err)
	}
	return nil
}

// IndexAge returns how old the cached index is, or an error if missing.
func IndexAge(root string) (time.Duration, error) {
	st, err := os.Stat(IndexPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("catalog not synced — run: wuji catalog sync")
		}
		return 0, err
	}
	return time.Since(st.ModTime()), nil
}
