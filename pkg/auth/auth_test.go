package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreSetGetRemove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Set("anthropic", Entry{Type: "api", Key: "sk-test"}); err != nil {
		t.Fatal(err)
	}
	entry, ok, err := store.Get("anthropic")
	if err != nil || !ok || entry.Key != "sk-test" {
		t.Fatalf("get: ok=%v entry=%+v err=%v", ok, entry, err)
	}
	if err := store.Remove("anthropic"); err != nil {
		t.Fatal(err)
	}
	_, ok, err = store.Get("anthropic")
	if err != nil || ok {
		t.Fatalf("expected removed, ok=%v err=%v", ok, err)
	}
}

func TestResolveKeyPrefersAuthStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	t.Setenv("WUJI_AUTH_PATH", path)
	t.Setenv("ANTHROPIC_API_KEY", "from-env")

	store, _ := Open(path)
	_ = store.Set("anthropic", Entry{Type: "api", Key: "from-auth"})

	got, err := ResolveKey("anthropic", "", "ANTHROPIC_API_KEY")
	if err != nil || got != "from-auth" {
		t.Fatalf("ResolveKey = %q err=%v", got, err)
	}
}

func TestResolveKeyEnvFallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	t.Setenv("WUJI_AUTH_PATH", path)
	t.Setenv("OPENAI_API_KEY", "env-key")

	got, err := ResolveKey("openai", "", "OPENAI_API_KEY")
	if err != nil || got != "env-key" {
		t.Fatalf("ResolveKey = %q err=%v", got, err)
	}
}

func TestResolveKeyMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	t.Setenv("WUJI_AUTH_PATH", path)
	_ = os.Unsetenv("ANTHROPIC_API_KEY")

	_, err := ResolveKey("anthropic", "", "ANTHROPIC_API_KEY")
	if err == nil {
		t.Fatal("expected error for missing credential")
	}
}

func TestDefaultPathUsesEnv(t *testing.T) {
	t.Setenv("WUJI_AUTH_PATH", "/tmp/custom-auth.json")
	if got := DefaultPath(); got != "/tmp/custom-auth.json" {
		t.Fatalf("DefaultPath() = %q", got)
	}
}

func TestResolveKeyInline(t *testing.T) {
	got, err := ResolveKey("openai", "inline-key", "OPENAI_API_KEY")
	if err != nil || got != "inline-key" {
		t.Fatalf("ResolveKey inline = %q err=%v", got, err)
	}
}

func TestResolveKeyFromEntryEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	t.Setenv("WUJI_AUTH_PATH", path)
	store, _ := Open(path)
	_ = store.Set("openai", Entry{Type: "api", Env: map[string]string{"OPENAI_API_KEY": "env-in-store"}})

	got, err := ResolveKey("openai", "", "OPENAI_API_KEY")
	if err != nil || got != "env-in-store" {
		t.Fatalf("ResolveKey = %q err=%v", got, err)
	}
}

func TestStoreCorruptJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.All(); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestDefaultPathXDG(t *testing.T) {
	t.Setenv("WUJI_AUTH_PATH", "")
	xdg := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_DATA_HOME", xdg)
	want := filepath.Join(xdg, "wuji", "auth.json")
	if got := DefaultPath(); got != want {
		t.Fatalf("DefaultPath() = %q want %q", got, want)
	}
}

func TestGetNormalizesProviderID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Set("openai/", Entry{Type: "api", Key: "sk"}); err != nil {
		t.Fatal(err)
	}
	entry, ok, err := store.Get("openai")
	if err != nil || !ok || entry.Key != "sk" {
		t.Fatalf("get normalized id: ok=%v entry=%+v err=%v", ok, entry, err)
	}
}

func TestResolveKeyFromFirstEnvValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	t.Setenv("WUJI_AUTH_PATH", path)
	store, _ := Open(path)
	_ = store.Set("demo", Entry{Type: "api", Env: map[string]string{"OTHER": "fallback-key"}})

	got, err := ResolveKey("demo", "", "MISSING_ENV")
	if err != nil || got != "fallback-key" {
		t.Fatalf("ResolveKey = %q err=%v", got, err)
	}
}

func TestStorePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if store.Path() != path {
		t.Fatalf("Path() = %q", store.Path())
	}
}
