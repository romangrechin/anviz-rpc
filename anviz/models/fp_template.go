package models

type FpTemplate struct {
	UserCode   uint64
	BackupCode uint8
	Data       []byte
}

func (ft *FpTemplate) Unmarshal(data []byte) error {
	ft.Data = data
	return nil
}

func (ft *FpTemplate) MarshalRequest() []byte {
	data := make([]byte, 6)
	bufCode := fromUInt64(ft.UserCode)
	copy(data[0:5], bufCode)
	data[5] = ft.BackupCode
	return data
}

func (ft *FpTemplate) Marshal() []byte {
	dataLength := 344
	if len(ft.Data) > 338 { // Iris
		dataLength = 1286
	}

	data := make([]byte, dataLength)
	bufCode := fromUInt64(ft.UserCode)
	copy(data[0:5], bufCode)
	data[5] = ft.BackupCode
	copy(data[6:], ft.Data)
	return data
}
