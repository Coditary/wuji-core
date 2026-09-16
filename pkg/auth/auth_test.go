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
