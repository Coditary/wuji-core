package api

import "testing"

func TestExtractJSONPath(t *testing.T) {
	raw := []byte(`{"data":[{"url":"https://example.com/x.png"}]}`)
	v, err := extractJSONPath(raw, "data[0].url")
	if err != nil {
		t.Fatal(err)
	}
	if v != "https://example.com/x.png" {
		t.Fatalf("got %q", v)
	}
}
