package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

func main() {
	avzScanner, err := syscall.LoadLibrary("AvzScanner.dll")
	if err != nil {
		log.Println("Error when load AvzScanner.dll: ", err)
		return
	}
	defer syscall.FreeLibrary(avzScanner)

	avzFindDevice, err := syscall.GetProcAddress(avzScanner, "AvzFindDevice")
	if err != nil {
		log.Println("Error when get proc address for AvzFindDevice: ", err)
		return
	}

	var names [8][128]byte

	r1, _, err := syscall.Syscall(avzFindDevice, 1, uintptr(unsafe.Pointer(&names)), 0, 0)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Println("Error when call function AvzFindDevice: ", err)
		return
	}

	ret := uint16(r1)
	fmt.Println("Function AvzFindDevice return code: ", ret)
	if ret != 1 {
		return
	}

	fmt.Println("Device names:")

	for _, v := range names {
		fmt.Println(string(v[:]))
	}
}
