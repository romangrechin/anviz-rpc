package device

import (
	"encoding/binary"
	"github.com/romangrechin/anviz-rpc/anviz/errors"
	"github.com/sigurn/crc16"
)

func marshal(cmd uint8, deviceId uint32, data []byte) []byte {
	buf := make([]byte, 10+len(data))
	buf[0] = 0xA5
	binary.BigEndian.PutUint32(buf[1:5], deviceId)
	buf[5] = cmd
	binary.BigEndian.PutUint16(buf[6:8], uint16(len(data)))
	if len(data) > 0 {
		copy(buf[8:8+len(data)], data)
	}
	sign(buf)
	return buf
}

func unmarshal(cmd uint8, buf []byte) (deviceId uint32, data []byte, err error) {
	if len(buf) < 11 || buf[0] != 0xA5 {
		err = errors.ErrInvalidResponseData
		return
	}

	if !isValid(buf) {
		err = errors.ErrInvalidResponseChecksum
		return
	}

	deviceId = binary.BigEndian.Uint32(buf[1:5])

	if buf[5] != cmd+0x80 {
		err = errors.ErrInvalidResponseAck
		return
	}

	if buf[6] != 0x0 {
		err = parseRet(buf[6])
		return
	}

	length := binary.BigEndian.Uint16(buf[7:9])
	if length > 0 {
		data = make([]byte, length)
		copy(data, buf[9:9+length])
	}

	return
}

func calcCrc(buf []byte) uint16 {
	return crc16.Checksum(buf, crc16.MakeTable(crc16.CRC16_MCRF4XX))
}

func sign(buf []byte) {
	crc := calcCrc(buf[:len(buf)-2])
	binary.LittleEndian.PutUint16(buf[len(buf)-2:], crc)
}

func isValid(buf []byte) bool {
	crc := calcCrc(buf[:len(buf)-2])
	return crc == binary.LittleEndian.Uint16(buf[len(buf)-2:])
}

func parseRet(ret uint8) error {
	return errors.ErrDeviceOperationFailed
}
