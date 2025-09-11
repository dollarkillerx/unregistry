package progress

import (
	"io"

	"github.com/schollz/progressbar/v3"
)

type Reader struct {
	io.Reader
	bar *progressbar.ProgressBar
}

func NewReader(r io.Reader, total int64, description string) *Reader {
	bar := progressbar.NewOptions64(
		total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	
	return &Reader{
		Reader: r,
		bar:    bar,
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		r.bar.Add(n)
	}
	if err == io.EOF {
		r.bar.Finish()
	}
	return n, err
}

func (r *Reader) Close() error {
	r.bar.Finish()
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type Writer struct {
	io.Writer
	bar *progressbar.ProgressBar
}

func NewWriter(w io.Writer, total int64, description string) *Writer {
	bar := progressbar.NewOptions64(
		total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	
	return &Writer{
		Writer: w,
		bar:    bar,
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if n > 0 {
		w.bar.Add(n)
	}
	return n, err
}

func (w *Writer) Close() error {
	w.bar.Finish()
	if closer, ok := w.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}