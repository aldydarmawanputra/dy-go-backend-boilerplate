package imageproc

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"strings"

	xdraw "golang.org/x/image/draw"
)

// maxPixels caps decoded image area to guard against decompression bombs
// (small file, enormous dimensions) that would exhaust CPU/RAM.
const maxPixels = 24_000_000 // ~24 megapixels

var ErrImageTooLarge = errors.New("image dimensions exceed the allowed limit")

func IsImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/jpeg") || strings.HasPrefix(contentType, "image/png")
}

// Thumbnail decodes an image and scales it to fit within maxW x maxH (preserving
// aspect ratio, never upscaling), re-encoding as JPEG. It rejects images whose
// declared dimensions exceed maxPixels before doing the expensive full decode.
func Thumbnail(data []byte, maxW, maxH int) ([]byte, string, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxPixels {
		return nil, "", ErrImageTooLarge
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	scale := min(float64(maxW)/float64(w), float64(maxH)/float64(h))
	if scale >= 1 {
		scale = 1
	}
	tw, th := int(float64(w)*scale), int(float64(h)*scale)
	if tw < 1 {
		tw = 1
	}
	if th < 1 {
		th = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, tw, th))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, b, xdraw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85}); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "image/jpeg", nil
}
