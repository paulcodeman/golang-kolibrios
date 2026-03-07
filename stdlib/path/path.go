package path

func Clean(name string) string {
	if name == "" {
		return "."
	}

	rooted := name[0] == '/'
	parts := make([]string, 0, 8)
	index := 0

	for index < len(name) {
		for index < len(name) && name[index] == '/' {
			index++
		}
		if index >= len(name) {
			break
		}

		next := index
		for next < len(name) && name[next] != '/' {
			next++
		}

		part := name[index:next]
		switch part {
		case "", ".":
		case "..":
			if rooted {
				if len(parts) > 0 {
					parts = parts[:len(parts)-1]
				}
			} else if len(parts) > 0 && parts[len(parts)-1] != ".." {
				parts = parts[:len(parts)-1]
			} else {
				parts = append(parts, part)
			}
		default:
			parts = append(parts, part)
		}

		index = next
	}

	if rooted {
		if len(parts) == 0 {
			return "/"
		}

		cleaned := "/" + parts[0]
		for index = 1; index < len(parts); index++ {
			cleaned += "/" + parts[index]
		}
		return cleaned
	}

	if len(parts) == 0 {
		return "."
	}

	cleaned := parts[0]
	for index = 1; index < len(parts); index++ {
		cleaned += "/" + parts[index]
	}
	return cleaned
}

func Join(elem ...string) string {
	joined := ""

	for index := 0; index < len(elem); index++ {
		part := elem[index]
		if part == "" {
			continue
		}

		if joined == "" {
			joined = part
			continue
		}

		if joined[len(joined)-1] == '/' {
			joined += part
			continue
		}

		joined += "/" + part
	}

	if joined == "" {
		return ""
	}

	return Clean(joined)
}

func Base(name string) string {
	if name == "" {
		return "."
	}

	cleaned := Clean(name)
	if cleaned == "/" {
		return "/"
	}

	slash := lastIndexByte(cleaned, '/')
	if slash < 0 {
		return cleaned
	}

	return cleaned[slash+1:]
}

func Dir(name string) string {
	if name == "" {
		return "."
	}

	cleaned := Clean(name)
	if cleaned == "/" {
		return "/"
	}

	slash := lastIndexByte(cleaned, '/')
	if slash < 0 {
		return "."
	}
	if slash == 0 {
		return "/"
	}

	return cleaned[:slash]
}

func Split(name string) (dir string, file string) {
	slash := lastIndexByte(name, '/')
	if slash < 0 {
		return "", name
	}

	return name[:slash+1], name[slash+1:]
}

func Ext(name string) string {
	base := Base(name)
	dot := lastIndexByte(base, '.')
	if dot < 0 {
		return ""
	}

	return base[dot:]
}

func IsAbs(name string) bool {
	return len(name) > 0 && name[0] == '/'
}

func lastIndexByte(name string, target byte) int {
	for index := len(name) - 1; index >= 0; index-- {
		if name[index] == target {
			return index
		}
	}

	return -1
}
