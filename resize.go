package main

import (
	"bytes"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/nfnt/resize"
	"image/png"
)

func resizeImage(img []byte, imageParams *ImageParams) ([]byte, error) {
	decoded, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return []byte{}, err
	}

	newImage := resize.Resize(uint(imageParams.Width), uint(imageParams.Height), decoded, resize.NearestNeighbor)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, newImage)
	return buf.Bytes(), err
}
