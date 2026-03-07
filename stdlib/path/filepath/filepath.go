package filepath

import (
	"os"
	"path"
)

const (
	Separator     = os.PathSeparator
	ListSeparator = os.PathListSeparator
)

func Clean(name string) string {
	return path.Clean(ToSlash(name))
}

func Join(elem ...string) string {
	normalized := make([]string, len(elem))
	for index := 0; index < len(elem); index++ {
		normalized[index] = ToSlash(elem[index])
	}

	return path.Join(normalized...)
}

func Base(name string) string {
	return path.Base(ToSlash(name))
}

func Dir(name string) string {
	return path.Dir(ToSlash(name))
}

func Split(name string) (dir string, file string) {
	return path.Split(ToSlash(name))
}

func Ext(name string) string {
	return path.Ext(ToSlash(name))
}

func IsAbs(name string) bool {
	return path.IsAbs(ToSlash(name))
}

func Abs(name string) (string, error) {
	if IsAbs(name) {
		return Clean(name), nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return Clean(path.Join(wd, ToSlash(name))), nil
}

func ToSlash(name string) string {
	normalized := make([]byte, len(name))
	for index := 0; index < len(name); index++ {
		if name[index] == '\\' {
			normalized[index] = '/'
			continue
		}

		normalized[index] = name[index]
	}

	return string(normalized)
}

func FromSlash(name string) string {
	if Separator == '/' {
		return name
	}

	normalized := make([]byte, len(name))
	for index := 0; index < len(name); index++ {
		if name[index] == '/' {
			normalized[index] = Separator
			continue
		}

		normalized[index] = name[index]
	}

	return string(normalized)
}

func VolumeName(name string) string {
	return ""
}
