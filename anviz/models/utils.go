package models

import (
	"encoding/binary"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func toInt(data []byte) int32 {
	buf := make([]byte, 3)
	copy(buf, data)

	return (int32(buf[0])&0xff)<<16 | (int32(buf[1])&0xff)<<8 | int32(buf[2])&0xff
}

func toInt64(data []byte) uint64 {
	buf := make([]byte, 8)
	copy(buf, data)

	return (uint64(buf[0])&0xff)<<32 | (uint64(buf[1])&0xff)<<24 | uint64(buf[2])&0xff<<16 |
		uint64(buf[3])&0xff<<8 | uint64(buf[4])&0xff
}

func fromInt(val int32) []byte {
	buf := make([]byte, 3)
	buf[0] = uint8((0xff0000 & val) >> 16)
	buf[1] = uint8((0x00ff00 & val) >> 8)
	buf[2] = uint8(0x0000ff & val)
	return buf
}

func fromUInt64(val uint64) []byte {
	buf := make([]byte, 5)
	buf[0] = uint8((0xff00000000 & val) >> 32)
	buf[1] = uint8((0x00ff000000 & val) >> 24)
	buf[2] = uint8((0x0000ff0000 & val) >> 16)
	buf[3] = uint8((0x000000ff00 & val) >> 8)
	buf[4] = uint8(0x00000000ff & val)
	return buf
}

func unpackDevicePassword(buf []byte) (uint32, uint8) {
	length := buf[0] >> 4
	password := (uint32(buf[0])&0x0f)<<16 | (uint32(buf[1])&0xff)<<8 | uint32(buf[2])&0xff
	return password, length
}

func packDevicePassword(password uint32, length uint8) []byte {
	buf := make([]byte, 3)
	buf[0] = uint8((0x0f0000&password)>>16) | ((0x0f & length) << 4)
	buf[1] = uint8((0x00ff00 & password) >> 8)
	buf[2] = uint8(0x0000ff & password)
	return buf
}

func unpackUserPassword(buf []byte, add byte) (uint32, uint8) {
	length := buf[0] >> 4
	// uint32(add & 0xff) << 20 |
	password := (uint32(buf[0])&0x0f)<<16 | (uint32(buf[1])&0xff)<<8 | uint32(buf[2])&0xff
	return password, length
}

func packUserPassword(password uint32, length uint8) ([]byte, byte) {
	buf := make([]byte, 3)
	buf[0] = uint8((0x000f0000&password)>>16) | ((0x0f & length) << 4)
	buf[1] = uint8((0x0000ff00 & password) >> 8)
	buf[2] = uint8(0x000000ff & password)
	return buf, uint8((0x0ff00000 & password) >> 20)
}

func isEmpty(data []byte) bool {
	for _, b := range data {
		if b != 0xff {
			return false
		}
	}
	return true
}

func utfToString(b []byte) string {
	var buf []byte
	for i := range b {
		if b[i] == 0x0 {
			break
		}

		buf = append(buf, b[i])
	}

	return string(buf)
}

func unicodeToString(b []byte) string {
	o := binary.BigEndian
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}

	var pos int
	for i := 0; i < len(utf); i++ {
		if utf[i] != 0 {
			pos = i
			break
		}
	}
	utf = utf[pos:]

	for i := len(utf) - 1; i > 0; i-- {
		if utf[i] != 0 {
			pos = i
			break
		}
	}
	utf = utf[:pos+1]

	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return strings.TrimSpace(string(utf16.Decode(utf)))
}

func stringToUnicode(val string, max int) []byte {
	o := binary.BigEndian
	buf := make([]byte, max)

	runes := utf16.Encode([]rune(val))
	for i := 0; i < len(runes); i++ {
		if len(buf)-1 < i*2 {
			break
		}
		o.PutUint16(buf[i*2:i*2+2], runes[i])
	}
	return buf
}

func stringToUtf8(val string, max int) []byte {
	buf := make([]byte, max)

	runeBuf := make([]byte, utf8.UTFMax)
	runes := []rune(val)
	length := 0
	for i := 0; i < len(runes); i++ {
		n := utf8.EncodeRune(runeBuf, runes[i])
		copy(buf[length:length+n], runeBuf[:n])
		length += n
		if length >= max {
			break
		}
	}
	return buf
}
