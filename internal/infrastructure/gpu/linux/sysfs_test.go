package linux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSysfsValue(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "value")

	err := os.WriteFile(
		path,
		[]byte("  hello world \n"),
		0644,
	)
	if err != nil {
		t.Fatalf("write test file: %v", err)
	}

	value, err := readSysfsValue(path)
	if err != nil {
		t.Fatalf("read sysfs value: %v", err)
	}

	if value != "hello world" {
		t.Fatalf(
			"expected %q, got %q",
			"hello world",
			value,
		)
	}
}

func TestReadSysfsUint64(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "value")

	err := os.WriteFile(
		path,
		[]byte("4294967296\n"),
		0644,
	)
	if err != nil {
		t.Fatalf("write test file: %v", err)
	}

	value, err := readSysfsUint64(path)
	if err != nil {
		t.Fatalf("read sysfs uint64: %v", err)
	}

	if value != 4294967296 {
		t.Fatalf(
			"expected %d, got %d",
			uint64(4294967296),
			value,
		)
	}
}