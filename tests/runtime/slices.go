package runtimeprobe

func SliceRoundtrip(src string) string {
	data := []byte(src)
	data = append(data, '!')

	dst := make([]byte, len(data))
	copy(dst, data)

	return string(dst)
}
