package imageproc

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestThumbnailDownscalesAndKeepsAspect(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 400, 200))
	for x := 0; x < 400; x++ {
		for y := 0; y < 200; y++ {
			src.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 100, A: 255})
		}
	}
	var in bytes.Buffer
	if err := png.Encode(&in, src); err != nil {
		t.Fatal(err)
	}

	out, ct, err := Thumbnail(in.Bytes(), 100, 100)
	if err != nil {
		t.Fatalf("Thumbnail: %v", err)
	}
	if ct != "image/jpeg" {
		t.Fatalf("content-type = %s", ct)
	}

	img, _, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode result: %v", err)
	}
	b := img.Bounds()
	if b.Dx() != 100 || b.Dy() != 50 {
		t.Fatalf("size = %dx%d, want 100x50 (aspect preserved)", b.Dx(), b.Dy())
	}
}

func TestIsImage(t *testing.T) {
	if !IsImage("image/png") || !IsImage("image/jpeg") {
		t.Fatal("png/jpeg should be images")
	}
	if IsImage("application/pdf") {
		t.Fatal("pdf should not be an image")
	}
}
