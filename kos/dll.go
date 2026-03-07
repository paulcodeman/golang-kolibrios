package kos

type DLLExportTable uint32

const ConsoleDLLPath = "/sys/lib/console.obj"

func LoadDLLFile(path string) DLLExportTable {
	return LoadDLLFileWithEncoding(path, EncodingUTF8)
}

func LoadDLLFileWithEncoding(path string, encoding StringEncoding) DLLExportTable {
	return DLLExportTable(LoadDLLWithEncoding(encoding, path))
}

func LoadDLLFileLegacy(path string) DLLExportTable {
	return DLLExportTable(LoadDLL(path))
}

func LoadConsoleDLL() DLLExportTable {
	return LoadDLLFile(ConsoleDLLPath)
}
