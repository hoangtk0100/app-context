package util

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func GetImageExtension(fileName string) string {
	fileExt := filepath.Ext(fileName)
	return strings.ReplaceAll(fileExt, ".", "")
}

func GetImageDimension(data []byte) (int, int, error) {
	fileBytes := bytes.NewBuffer(data)
	img, _, err := image.DecodeConfig(fileBytes)
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}

	return img.Width, img.Height, nil
}
