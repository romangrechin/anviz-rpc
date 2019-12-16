package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"log"
	"syscall"
	"unsafe"
)

func uintptrToBytes(u uintptr) []byte {
	return (*[2048]byte)(unsafe.Pointer(u))[:]
}

func main() {
	anvizComponent, err := syscall.LoadLibrary("AnvizComponent.dll")
	if err != nil {
		log.Println("Error when load AnvizComponent.dll: ", err)
		return
	}
	defer syscall.FreeLibrary(anvizComponent)

	getClassNames, err := syscall.GetProcAddress(anvizComponent, "GetClassNames")
	if err != nil {
		log.Println("Error when get proc address for GetClassNames: ", err)
		return
	}

	r1, _, err := syscall.Syscall(getClassNames, 0, 0, 0, 0)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Println("Error when call function AvzFindDevice: ", err)
		return
	}

	b := uintptrToBytes(r1)

	dec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	out, err := dec.Bytes(b)

	if err != nil {
		log.Println("Convert to UTF16 error: ", err)
	}

	i := bytes.IndexByte(out, 0)
	if i == -1 {
		i = len(out)
	}

	fmt.Println(string(out[:i]))
}
