package web

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomic_ReplacesContentAtomically(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "status.json")

	if err := writeFileAtomic(path, []byte("{\"zones\":[{\"id\":\"old\"}]}"), 0644); err != nil {
		t.Fatalf("first write: %v", err)
	}

	if err := writeFileAtomic(path, []byte("{\"zones\":[{\"id\":\"new\"}]}"), 0644); err != nil {
		t.Fatalf("second write: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	if got := string(data); got != "{\"zones\":[{\"id\":\"new\"}]}" {
		t.Fatalf("unexpected content: %q", got)
	}
}
