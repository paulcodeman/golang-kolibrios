package kos

func FreeRAMKB() uint32 {
	return GetFreeRAM()
}

func TotalRAMKB() uint32 {
	return GetTotalRAM()
}
