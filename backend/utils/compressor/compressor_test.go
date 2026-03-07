package compressor

import (
	"os"
	"testing"
)

func TestCompressImage_PNG(t *testing.T) {
	f, err := os.Open("../../test/test.png")
	if err != nil {
		t.Fatalf("failed to open test.png: %v", err)
	}
	defer f.Close()

	out, ext, err := CompressImage(f)
	if err != nil {
		t.Fatalf("CompressImage failed: %v", err)
	}
	defer os.Remove(out.Name())
	defer out.Close()

	if ext != "webp" {
		t.Errorf("expected ext \"webp\", got %q", ext)
	}

	info, err := out.Stat()
	if err != nil {
		t.Fatalf("failed to stat output: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestCompressImage_JPEG(t *testing.T) {
	f, err := os.Open("../../test/test.jpg")
	if err != nil {
		t.Fatalf("failed to open test.jpg: %v", err)
	}
	defer f.Close()

	out, ext, err := CompressImage(f)
	if err != nil {
		t.Fatalf("CompressImage failed: %v", err)
	}
	defer os.Remove(out.Name())
	defer out.Close()

	if ext != "webp" {
		t.Errorf("expected ext \"webp\", got %q", ext)
	}

	info, err := out.Stat()
	if err != nil {
		t.Fatalf("failed to stat output: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}
