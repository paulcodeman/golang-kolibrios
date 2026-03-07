package os

import (
	"io"
	"kos"
	"syscall"
)

type FileMode uint32

const (
	ModeDir FileMode = 1 << 31
)

const (
	O_RDONLY int = 0
	O_WRONLY int = 1
	O_RDWR   int = 2

	O_CREATE int = 0x40
	O_TRUNC  int = 0x200
	O_APPEND int = 0x400
)

type osError struct {
	text string
}

func (err *osError) Error() string {
	return err.text
}

var ErrInvalid = &osError{text: "invalid argument"}
var ErrPermission = &osError{text: "permission denied"}
var ErrExist = &osError{text: "file already exists"}
var ErrNotExist = &osError{text: "file does not exist"}
var ErrClosed = &osError{text: "file already closed"}

var Stdin = newDescriptorFile("stdin", int(kos.StdinFD), true, false)
var Stdout = newDescriptorFile("stdout", int(kos.StdoutFD), false, true)
var Stderr = newDescriptorFile("stderr", int(kos.StderrFD), false, true)

type PathError struct {
	Op   string
	Path string
	Err  error
}

func (err *PathError) Error() string {
	if err == nil {
		return ""
	}
	if err.Err == nil {
		return err.Op + " " + err.Path
	}

	return err.Op + " " + err.Path + ": " + err.Err.Error()
}

func (err *PathError) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Err
}

type LinkError struct {
	Op  string
	Old string
	New string
	Err error
}

func (err *LinkError) Error() string {
	if err == nil {
		return ""
	}
	if err.Err == nil {
		return err.Op + " " + err.Old + " " + err.New
	}

	return err.Op + " " + err.Old + " " + err.New + ": " + err.Err.Error()
}

func (err *LinkError) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Err
}

type statusError struct {
	status kos.FileSystemStatus
	text   string
}

func (err *statusError) Error() string {
	return err.text
}

type File struct {
	name     string
	fd       int
	offset   uint64
	readable bool
	writable bool
	append   bool
	closed   bool
	fdBacked bool
}

func Getwd() (dir string, err error) {
	dir = kos.CurrentFolder()
	if dir == "" {
		return "", &PathError{Op: "getwd", Path: "", Err: ErrInvalid}
	}

	return dir, nil
}

func ReadFile(name string) ([]byte, error) {
	data, status := kos.ReadAllFile(name)
	if status == kos.FileSystemOK || status == kos.FileSystemEOF {
		return data, nil
	}

	return nil, wrapPathError("read", name, status)
}

func WriteFile(name string, data []byte, perm FileMode) error {
	written, status := kos.CreateOrRewriteFile(name, data)
	if status != kos.FileSystemOK {
		return wrapPathError("write", name, status)
	}
	if int(written) != len(data) {
		return &PathError{Op: "write", Path: name, Err: io.ErrShortWrite}
	}

	return nil
}

func Mkdir(name string, perm FileMode) error {
	status := kos.CreateDirectory(name)
	if status != kos.FileSystemOK {
		return wrapPathError("mkdir", name, status)
	}

	return nil
}

func Remove(name string) error {
	status := kos.DeletePath(name)
	if status != kos.FileSystemOK {
		return wrapPathError("remove", name, status)
	}

	return nil
}

func Rename(oldpath string, newpath string) error {
	status := kos.RenamePath(oldpath, newpath)
	if status != kos.FileSystemOK {
		return &LinkError{
			Op:  "rename",
			Old: oldpath,
			New: newpath,
			Err: statusToError(status),
		}
	}

	return nil
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func Create(name string) (*File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0)
}

func Pipe() (reader *File, writer *File, err error) {
	var pipefd [2]int

	if err = syscall.Pipe(pipefd[:]); err != nil {
		return nil, nil, err
	}

	reader = newDescriptorFile("pipe[0]", pipefd[0], true, false)
	writer = newDescriptorFile("pipe[1]", pipefd[1], false, true)
	return reader, writer, nil
}

func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	accessMode := flag & 3
	readable := accessMode == O_RDONLY || accessMode == O_RDWR
	writable := accessMode == O_WRONLY || accessMode == O_RDWR

	if flag&O_TRUNC != 0 && !writable {
		return nil, &PathError{Op: "open", Path: name, Err: ErrInvalid}
	}
	if flag&O_APPEND != 0 && !writable {
		return nil, &PathError{Op: "open", Path: name, Err: ErrInvalid}
	}

	if flag&O_CREATE != 0 {
		_, status := kos.GetPathInfo(name)
		if status == kos.FileSystemNotFound {
			_, status = kos.CreateOrRewriteFile(name, nil)
			if status != kos.FileSystemOK {
				return nil, wrapPathError("open", name, status)
			}
		} else if status != kos.FileSystemOK {
			return nil, wrapPathError("open", name, status)
		}
	}

	if flag&O_TRUNC != 0 {
		_, status := kos.CreateOrRewriteFile(name, nil)
		if status != kos.FileSystemOK {
			return nil, wrapPathError("open", name, status)
		}
	}

	info, status := kos.GetPathInfo(name)
	if status != kos.FileSystemOK {
		return nil, wrapPathError("open", name, status)
	}

	file := &File{
		name:     name,
		readable: readable,
		writable: writable,
		append:   flag&O_APPEND != 0,
	}
	if file.append {
		file.offset = info.Size
	}

	return file, nil
}

func (file *File) Name() string {
	if file == nil {
		return ""
	}

	return file.name
}

func (file *File) Close() error {
	if file == nil {
		return &PathError{Op: "close", Path: "", Err: ErrInvalid}
	}
	if file.closed {
		return &PathError{Op: "close", Path: file.name, Err: ErrClosed}
	}

	file.closed = true
	return nil
}

func (file *File) Read(buffer []byte) (int, error) {
	if err := file.ensureReadable("read"); err != nil {
		return 0, err
	}
	if len(buffer) == 0 {
		return 0, nil
	}
	if file.fdBacked {
		read, err := syscall.Read(file.fd, buffer)
		if err != nil {
			return read, &PathError{Op: "read", Path: file.name, Err: err}
		}
		if read == 0 {
			return 0, io.EOF
		}
		return read, nil
	}

	read, status := kos.ReadFile(file.name, buffer, file.offset)
	file.offset += uint64(read)

	switch status {
	case kos.FileSystemOK:
		if read == 0 {
			return 0, io.EOF
		}
		return int(read), nil
	case kos.FileSystemEOF:
		if read > 0 {
			return int(read), io.EOF
		}
		return 0, io.EOF
	default:
		return int(read), wrapPathError("read", file.name, status)
	}
}

func (file *File) Write(buffer []byte) (int, error) {
	if err := file.ensureWritable("write"); err != nil {
		return 0, err
	}
	if len(buffer) == 0 {
		return 0, nil
	}
	if file.fdBacked {
		written, err := syscall.Write(file.fd, buffer)
		if err != nil {
			return written, &PathError{Op: "write", Path: file.name, Err: err}
		}
		if written != len(buffer) {
			return written, io.ErrShortWrite
		}
		return written, nil
	}

	if file.append {
		info, status := kos.GetPathInfo(file.name)
		if status != kos.FileSystemOK {
			return 0, wrapPathError("write", file.name, status)
		}
		file.offset = info.Size
	}

	written, status := kos.WriteFile(file.name, buffer, file.offset)
	file.offset += uint64(written)

	if status != kos.FileSystemOK {
		return int(written), wrapPathError("write", file.name, status)
	}
	if int(written) != len(buffer) {
		return int(written), io.ErrShortWrite
	}

	return int(written), nil
}

func (file *File) ensureReadable(op string) error {
	if file == nil {
		return &PathError{Op: op, Path: "", Err: ErrInvalid}
	}
	if file.closed {
		return &PathError{Op: op, Path: file.name, Err: ErrClosed}
	}
	if !file.readable {
		return &PathError{Op: op, Path: file.name, Err: ErrPermission}
	}

	return nil
}

func (file *File) ensureWritable(op string) error {
	if file == nil {
		return &PathError{Op: op, Path: "", Err: ErrInvalid}
	}
	if file.closed {
		return &PathError{Op: op, Path: file.name, Err: ErrClosed}
	}
	if !file.writable {
		return &PathError{Op: op, Path: file.name, Err: ErrPermission}
	}

	return nil
}

func wrapPathError(op string, name string, status kos.FileSystemStatus) error {
	return &PathError{
		Op:   op,
		Path: name,
		Err:  statusToError(status),
	}
}

func statusToError(status kos.FileSystemStatus) error {
	switch status {
	case kos.FileSystemOK:
		return nil
	case kos.FileSystemNotFound:
		return ErrNotExist
	case kos.FileSystemAccessDenied:
		return ErrPermission
	case kos.FileSystemUnsupported, kos.FileSystemBadPointer:
		return ErrInvalid
	case kos.FileSystemDiskFull:
		return &statusError{status: status, text: "disk full"}
	case kos.FileSystemInternalError:
		return &statusError{status: status, text: "internal error"}
	case kos.FileSystemDeviceError:
		return &statusError{status: status, text: "device error"}
	case kos.FileSystemNeedsMoreMemory:
		return &statusError{status: status, text: "not enough memory"}
	case kos.FileSystemEOF:
		return io.EOF
	}

	return &statusError{
		status: status,
		text:   "filesystem status " + formatStatus(status),
	}
}

var decimalDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func formatStatus(status kos.FileSystemStatus) string {
	return formatUint32(uint32(status))
}

func formatUint32(value uint32) string {
	if value < 10 {
		return decimalDigits[value]
	}

	return formatUint32(value/10) + decimalDigits[value%10]
}

func newDescriptorFile(name string, fd int, readable bool, writable bool) *File {
	return &File{
		name:     name,
		fd:       fd,
		readable: readable,
		writable: writable,
		fdBacked: true,
	}
}
