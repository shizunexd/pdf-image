package utils

import (
	"bytes"
	"errors"
	"image/jpeg"
	"io"

	"github.com/gen2brain/go-fitz"
)

func Convert(fileReader io.Reader, pageNum int) ([]byte, error) {
	doc, err := fitz.NewFromReader(fileReader)
	if err != nil {
		return nil, err
	}

	defer doc.Close()

	if err != nil {
		return nil, err
	}

	// Check if page is in range
	if pageNum > doc.NumPage()-1 || pageNum < -doc.NumPage() {
		return nil, errors.New("Page out of range")
	}

	// Convert negative index
	if pageNum < 0 {
		pageNum = doc.NumPage() + pageNum
	}

	// Extract pages as images
	img, err := doc.Image(pageNum)
	if err != nil {
		return nil, err
	}
	var imgBytes bytes.Buffer
	err = jpeg.Encode(&imgBytes, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return nil, err
	}
	return imgBytes.Bytes(), nil
}
