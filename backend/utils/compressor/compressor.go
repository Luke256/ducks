package compressor

import (
	"fmt"
	"image"
	"io"
	"os"

	"github.com/gen2brain/webp"
)

var ErrInvalidImage = fmt.Errorf("invalid image format")

// CompressImage 画像を圧縮し、webp形式でファイルを返します
func CompressImage(src io.ReadSeekCloser) (*os.File, string, error) {
	srcImage, _, err := image.Decode(src)
	if err != nil {
		if err == image.ErrFormat {
			return nil, "", ErrInvalidImage
		}
		return nil, "", err
	}

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, "", err
	}

	options := webp.Options{
		Lossless: false,
		Quality:  75,
	}
	if err := webp.Encode(tmpFile, srcImage, options); err != nil {
		return nil, "", err
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return nil, "", err
	}

	return tmpFile, "webp", nil
}