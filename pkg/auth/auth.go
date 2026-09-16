package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Entry is one stored provider credential (OpenCode-compatible shape).
type Entry struct {
	Type string            `json:"type"` // api | oauth (oauth reserved)
	Key  string            `json:"key,omitempty"`
	Env  map[string]string `json:"env,omitempty"`
}

// Store persists provider credentials in auth.json.
type Store struct {
	path string
	mu   sync.RWMutex
}

// DefaultPath returns the global credential file (~/.local/share/wuji/auth.json).
func DefaultPath() string {
	if v := strings.TrimSpace(os.Getenv("WUJI_AUTH_PATH")); v != "" {
		return v
	}
	base := strings.TrimSpace(os.Getenv("XDG_DATA_HOME"))
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(".wuji", "auth.json")
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "wuji", "auth.json")
}

// Open loads or creates a credential store at path (empty = DefaultPath).
func Open(path string) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultPath()
	}
	return &Store{path: path}, nil
}

// Path returns the auth file location.
func (s *Store) Path() string { return s.path }

// All returns all stored credentials.
func (s *Store) All() (map[string]Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readLocked()
}

func (s *Store) readLocked() (map[string]Entry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]Entry{}, nil
		}
		return nil, fmt.Errorf("read auth store: %w", err)
	}
	var out map[string]Entry
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("parse auth store: %w", err)
	}
	if out == nil {
		out = map[string]Entry{}
	}
	return out, nil
}

// Get returns a credential entry for providerID.
func (s *Store) Get(providerID string) (Entry, bool, error) {
	all, err := s.All()
	if err != nil {
		return Entry{}, false, err
	}
	entry, ok := all[normalizeID(providerID)]
	return entry, ok, nil
}

// Set stores a credential for providerID.
func (s *Store) Set(providerID string, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.readLocked()
	if err != nil {
		return err
	}
	if all == nil {
		all = map[string]Entry{}
	}
	all[normalizeID(providerID)] = entry
	return s.writeLocked(all)
}

// Remove deletes a provider credential.
func (s *Store) Remove(providerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.readLocked()
	if err != nil {
		return err
	}
	delete(all, normalizeID(providerID))
	return s.writeLocked(all)
}

func (s *Store) writeLocked(all map[string]Entry) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create auth dir: %w", err)
	}
	data, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return fmt.Errorf("encode auth store: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("write auth store: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename auth store: %w", err)
	}
	return nil
}

// ResolveKey returns a credential for providerID and env var name.
// Priority: inline config key, auth store (Key or Env), process environment.
func ResolveKey(providerID, inline, envName string) (string, error) {
	if strings.TrimSpace(inline) != "" {
		return inline, nil
	}
	envName = strings.TrimSpace(envName)
	store, err := Open("")
	if err == nil {
		if entry, ok, _ := store.Get(providerID); ok {
			if v := strings.TrimSpace(entry.Key); v != "" {
				return v, nil
			}
			if entry.Env != nil {
				if v := strings.TrimSpace(entry.Env[envName]); v != "" {
					return v, nil
				}
				for _, v := range entry.Env {
					if strings.TrimSpace(v) != "" {
						return strings.TrimSpace(v), nil
					}
				}
			}
		}
	}
	if envName == "" {
		return "", nil
	}
	v := strings.TrimSpace(os.Getenv(envName))
	if v == "" {
		return "", fmt.Errorf("credential for %q is not set — run: wuji provider login %s", providerID, providerID)
	}
	return v, nil
}

func normalizeID(id string) string {
	return strings.TrimRight(strings.TrimSpace(id), "/")
}
