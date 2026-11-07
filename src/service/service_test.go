package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBinaryCopyAndServiceIdentity(t *testing.T) {
	if GetBinaryPath() == "" {
		t.Fatal("binary path is empty")
	}
	if manager := DetectServiceManager(); manager < ServiceUnknown || manager > ServiceBSDRC {
		t.Fatalf("invalid service manager: %d", manager)
	}
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "nested", "binary")
	if err := os.WriteFile(source, []byte("binary contents"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := copyBinary(source, destination); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(destination)
	if err != nil || string(got) != "binary contents" {
		t.Fatalf("copied data = %q, %v", got, err)
	}
	if err := copyBinary(filepath.Join(root, "missing"), filepath.Join(root, "out")); err == nil {
		t.Fatal("copyBinary accepted missing source")
	}
}
