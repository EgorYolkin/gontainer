package reader

import (
	"bufio"
	"io"
)

type SPReader struct {
	reader *bufio.Reader
}

func NewSPReader(rd io.Reader) *SPReader {
	return &SPReader{reader: bufio.NewReader(rd)}
}
