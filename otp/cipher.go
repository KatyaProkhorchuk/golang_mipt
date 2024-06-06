//go:build !solution

package otp

import (
	"io"
)

type xorReader struct {
	r    io.Reader
	prng io.Reader
}

func (x xorReader) Read(p []byte) (n int, err error) {
	n, err = x.r.Read(p)
	prngBuffer := make([]byte, n)
	_, _ = x.prng.Read(prngBuffer)
	for i := 0; i < n; i++ {
		p[i] = p[i] ^ prngBuffer[i]
	}
	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &xorReader{r: r, prng: prng}
}

type xorWriter struct {
	w    io.Writer
	prng io.Reader
}

func (x xorWriter) Write(p []byte) (n int, err error) {
	prngBuffer := make([]byte, len(p))
	n, err = x.prng.Read(prngBuffer)
	if err != nil {
		return 0, err
	}
	for i := 0; i < n; i++ {
		prngBuffer[i] ^= p[i]
	}
	n, err = x.w.Write(prngBuffer)
	return n, err
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &xorWriter{w: w, prng: prng}
}
