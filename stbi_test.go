package stbi

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// testPNG encodes a small image with a gradient and partial alpha, to catch stride and premultiplication mistakes.
func testPNG(t *testing.T) (*image.NRGBA, []byte) {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, 7, 5))
	for y := range 5 {
		for x := range 7 {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x * 30), G: uint8(y * 50), B: 200, A: uint8(40 + x*y*5)})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}

	return img, buf.Bytes()
}

func TestLoadMemory(t *testing.T) {
	want, data := testPNG(t)

	got, err := LoadMemory(data)
	if err != nil {
		t.Fatal(err)
	}

	if got.Rect != want.Rect || got.Stride != want.Stride || !bytes.Equal(got.Pix, want.Pix) {
		t.Fatalf("decoded image differs from source")
	}

	if got.At(3, 2) != want.At(3, 2) {
		t.Fatalf("At(3, 2) = %v, want %v", got.At(3, 2), want.At(3, 2))
	}
}

func TestLoadMemoryRaw(t *testing.T) {
	want, data := testPNG(t)

	pix, width, height, err := LoadMemoryRaw(data)
	if err != nil {
		t.Fatal(err)
	}
	defer Free(pix)

	if width != 7 || height != 5 || !bytes.Equal(pix, want.Pix) {
		t.Fatalf("raw pixels differ from source")
	}
}

func TestLoadAndLoadf(t *testing.T) {
	want, data := testPNG(t)

	path := filepath.Join(t.TempDir(), "test.png")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got.Pix, want.Pix) {
		t.Fatalf("Load pixels differ from source")
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err = LoadFile(f)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got.Pix, want.Pix) {
		t.Fatalf("LoadFile pixels differ from source")
	}

	// stb_image converts LDR images to linear floats, so only alpha maps straight to [0, 1].
	fimg, err := Loadf(path)
	if err != nil {
		t.Fatal(err)
	}

	if fimg.Rect != want.Rect || fimg.Stride != 7*4 {
		t.Fatalf("Loadf rect/stride = %v/%d", fimg.Rect, fimg.Stride)
	}

	i := 2*fimg.Stride + 3*4 + 3
	if a := fimg.Pix[i]; a < float32(want.Pix[2*want.Stride+3*4+3])/255-0.01 || a > float32(want.Pix[2*want.Stride+3*4+3])/255+0.01 {
		t.Fatalf("Loadf alpha = %v", a)
	}

	fpix, width, height, err := LoadfRaw(path)
	if err != nil {
		t.Fatal(err)
	}
	defer Free(fpix)

	if width != 7 || height != 5 || len(fpix) != len(fimg.Pix) {
		t.Fatalf("LoadfRaw size = %dx%d, %d values", width, height, len(fpix))
	}
}

func TestErrors(t *testing.T) {
	if _, err := LoadMemory(nil); err == nil {
		t.Error("LoadMemory(nil) succeeded")
	}

	if _, err := LoadMemory([]byte("not an image")); err == nil {
		t.Error("LoadMemory(garbage) succeeded")
	}

	if _, err := Load(filepath.Join(t.TempDir(), "missing.png")); err == nil {
		t.Error("Load(missing) succeeded")
	}

	Free[byte](nil)
}
