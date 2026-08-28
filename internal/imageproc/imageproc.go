package imageproc

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"strings"

	xdraw "golang.org/x/image/draw"
)

func IsImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/jpeg") || strings.HasPrefix(contentType, "image/png")
}

// Thumbnail decodes an image and scales it to fit within maxW x maxH (preserving
// aspect ratio, never upscaling), re-encoding as JPEG.
func Thumbnail(r io.Reader, maxW, maxH int) ([]byte, string, error) {
	img, _, err := image.Decode(r)
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
