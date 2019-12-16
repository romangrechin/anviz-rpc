package anviz

type Command struct {
	Cmd      uint8
	DeviceId uint32
	Data     []byte
}

type Response struct {
	Ack      uint8
	Ret      uint8
	DeviceId uint32
	Data     []byte
}
