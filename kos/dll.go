package kos

type DLLExportTable uint32
type DLLProc uint32

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

func LookupDLLExport(table DLLExportTable, name string) DLLProc {
	namePtr, _ := stringAddress(name)
	if table == 0 || namePtr == nil {
		return 0
	}

	proc := DLLProc(LookupDLLExportRaw(uint32(table), namePtr))
	freeCString(namePtr)
	return proc
}

func (table DLLExportTable) Lookup(name string) DLLProc {
	return LookupDLLExport(table, name)
}

func (proc DLLProc) Valid() bool {
	return proc != 0
}
