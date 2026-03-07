package bufio

import (
	"bytes"
	"errors"
	"io"
)

const (
	defaultBufSize   = 4096
	MaxScanTokenSize = 64 * 1024
)

var ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
var ErrTooLong = errors.New("bufio.Scanner: token too long")
var errBadScanAdvance = errors.New("bufio.Scanner: bad split advance")
var errReaderNoProgress = errors.New("bufio.Scanner: reader returned no data")

type Reader struct {
	reader    io.Reader
	buffer    []byte
	lastByte  byte
	canUnread bool
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader: reader,
		buffer: make([]byte, 0, defaultBufSize),
	}
}

func (reader *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	reader.canUnread = false
	if len(reader.buffer) > 0 {
		n = copy(p, reader.buffer)
		reader.buffer = reader.buffer[n:]
		if n > 0 {
			reader.lastByte = p[n-1]
			return n, nil
		}
	}

	n, err = reader.reader.Read(p)
	if n > 0 {
		reader.lastByte = p[n-1]
	}
	return n, err
}

func (reader *Reader) ReadByte() (byte, error) {
	if len(reader.buffer) > 0 {
		value := reader.buffer[0]
		reader.buffer = reader.buffer[1:]
		reader.lastByte = value
		reader.canUnread = true
		return value, nil
	}

	var single [1]byte
	for {
		read, err := reader.reader.Read(single[:])
		if read > 0 {
			reader.lastByte = single[0]
			reader.canUnread = true
			return single[0], nil
		}
		if err != nil {
			return 0, err
		}
	}
}

func (reader *Reader) UnreadByte() error {
	if !reader.canUnread {
		return ErrInvalidUnreadByte
	}

	reader.buffer = append([]byte{reader.lastByte}, reader.buffer...)
	reader.canUnread = false
	return nil
}

func (reader *Reader) ReadBytes(delim byte) ([]byte, error) {
	data := make([]byte, 0, 64)

	for {
		value, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF && len(data) > 0 {
				return data, err
			}
			return data, err
		}

		data = append(data, value)
		if value == delim {
			return data, nil
		}
	}
}

func (reader *Reader) ReadString(delim byte) (string, error) {
	data, err := reader.ReadBytes(delim)
	return string(data), err
}

type Writer struct {
	writer io.Writer
	buffer []byte
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
		buffer: make([]byte, 0, defaultBufSize),
	}
}

func (writer *Writer) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if len(writer.buffer) == 0 && len(p) >= defaultBufSize {
			written, writeErr := writer.writer.Write(p)
			n += written
			if writeErr != nil {
				return n, writeErr
			}
			if written != len(p) {
				return n, io.ErrShortWrite
			}
			return n, nil
		}

		available := defaultBufSize - len(writer.buffer)
		if available == 0 {
			if err = writer.Flush(); err != nil {
				return n, err
			}
			continue
		}
		if available > len(p) {
			available = len(p)
		}

		writer.buffer = append(writer.buffer, p[:available]...)
		p = p[available:]
		n += available

		if len(writer.buffer) == defaultBufSize {
			if err = writer.Flush(); err != nil {
				return n, err
			}
		}
	}

	return n, nil
}

func (writer *Writer) WriteByte(value byte) error {
	_, err := writer.Write([]byte{value})
	return err
}

func (writer *Writer) WriteString(value string) (n int, err error) {
	return writer.Write([]byte(value))
}

func (writer *Writer) Flush() error {
	if len(writer.buffer) == 0 {
		return nil
	}

	written, err := writer.writer.Write(writer.buffer)
	if err != nil {
		writer.buffer = writer.buffer[written:]
		return err
	}
	if written != len(writer.buffer) {
		writer.buffer = writer.buffer[written:]
		return io.ErrShortWrite
	}

	writer.buffer = writer.buffer[:0]
	return nil
}

type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

type Scanner struct {
	reader       io.Reader
	split        SplitFunc
	buffer       []byte
	token        []byte
	maxTokenSize int
	done         bool
	readDone     bool
	readErr      error
	readScratch  [256]byte
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader:       reader,
		split:        ScanLines,
		maxTokenSize: MaxScanTokenSize,
	}
}

func (scanner *Scanner) Buffer(buffer []byte, max int) {
	if max > 0 {
		scanner.maxTokenSize = max
	}
	if buffer != nil {
		scanner.buffer = buffer[:0]
	}
}

func (scanner *Scanner) Split(split SplitFunc) {
	if split == nil {
		scanner.split = ScanLines
		return
	}

	scanner.split = split
}

func (scanner *Scanner) Scan() bool {
	if scanner.done {
		return false
	}

	for {
		atEOF := scanner.readDone
		advance, token, err := scanner.split(scanner.buffer, atEOF)
		if err != nil {
			scanner.done = true
			scanner.readErr = err
			return false
		}
		if advance < 0 || advance > len(scanner.buffer) {
			scanner.done = true
			scanner.readErr = errBadScanAdvance
			return false
		}
		if token != nil {
			scanner.buffer = scanner.buffer[advance:]
			scanner.token = append(scanner.token[:0], token...)
			return true
		}
		if atEOF {
			scanner.done = true
			if scanner.readErr == io.EOF {
				scanner.readErr = nil
			}
			return false
		}
		if advance > 0 {
			scanner.buffer = scanner.buffer[advance:]
			continue
		}
		if len(scanner.buffer) >= scanner.maxTokenSize {
			scanner.done = true
			scanner.readErr = ErrTooLong
			return false
		}

		read, err := scanner.reader.Read(scanner.readScratch[:])
		if read > 0 {
			scanner.buffer = append(scanner.buffer, scanner.readScratch[:read]...)
			if len(scanner.buffer) > scanner.maxTokenSize {
				scanner.done = true
				scanner.readErr = ErrTooLong
				return false
			}
		}
		if err != nil {
			scanner.readDone = true
			scanner.readErr = err
			continue
		}
		if read == 0 {
			scanner.done = true
			scanner.readErr = errReaderNoProgress
			return false
		}
	}
}

func (scanner *Scanner) Bytes() []byte {
	return scanner.token
}

func (scanner *Scanner) Text() string {
	return string(scanner.token)
}

func (scanner *Scanner) Err() error {
	return scanner.readErr
}

func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) > 0 {
		return 1, data[:1], nil
	}
	if atEOF {
		return 0, nil, nil
	}

	return 0, nil, nil
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if index := bytes.IndexByte(data, '\n'); index >= 0 {
		line := data[:index]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		return index + 1, line, nil
	}
	if atEOF && len(data) > 0 {
		if data[len(data)-1] == '\r' {
			return len(data), data[:len(data)-1], nil
		}
		return len(data), data, nil
	}

	return 0, nil, nil
}

func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for start < len(data) && isSpace(data[start]) {
		start++
	}
	if start >= len(data) {
		if atEOF {
			return len(data), nil, nil
		}
		return start, nil, nil
	}

	for index := start; index < len(data); index++ {
		if isSpace(data[index]) {
			return index + 1, data[start:index], nil
		}
	}
	if atEOF {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}

func isSpace(value byte) bool {
	switch value {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}

	return false
}
