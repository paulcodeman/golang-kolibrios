package io

type ioError struct {
	text string
}

func (err *ioError) Error() string {
	return err.text
}

var EOF = &ioError{text: "EOF"}
var ErrShortWrite = &ioError{text: "short write"}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

type ReadWriter interface {
	Reader
	Writer
}

type ReadCloser interface {
	Reader
	Closer
}

type WriteCloser interface {
	Writer
	Closer
}

type StringWriter interface {
	WriteString(s string) (n int, err error)
}

func ReadAll(r Reader) ([]byte, error) {
	data := make([]byte, 0, 512)
	buffer := make([]byte, 512)

	for {
		read, err := r.Read(buffer)
		if read > 0 {
			data = append(data, buffer[:read]...)
		}

		if err != nil {
			if err == EOF {
				return data, nil
			}

			return data, err
		}
	}
}

func Copy(dst Writer, src Reader) (written int64, err error) {
	return CopyBuffer(dst, src, nil)
}

func CopyBuffer(dst Writer, src Reader, buffer []byte) (written int64, err error) {
	if len(buffer) == 0 {
		buffer = make([]byte, 512)
	}

	for {
		read, readErr := src.Read(buffer)
		if read > 0 {
			wrote, writeErr := dst.Write(buffer[:read])
			written += int64(wrote)

			if writeErr != nil {
				return written, writeErr
			}
			if wrote != read {
				return written, ErrShortWrite
			}
		}

		if readErr != nil {
			if readErr == EOF {
				return written, nil
			}

			return written, readErr
		}
	}
}

func WriteString(w Writer, s string) (n int, err error) {
	return w.Write([]byte(s))
}
