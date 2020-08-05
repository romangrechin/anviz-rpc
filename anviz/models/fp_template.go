package models

import "log"

type FpTemplate struct {
	UserCode   uint64
	BackupCode uint8
	Data       []byte
}

func (ft *FpTemplate) Unmarshal(data []byte) error {
	log.Println("len data out: ", len(data))
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
	log.Println("len data in: ", len(ft.Data))
	data := make([]byte, 6+len(ft.Data))
	bufCode := fromUInt64(ft.UserCode)
	copy(data[0:5], bufCode)
	data[5] = ft.BackupCode
	copy(data[6:], ft.Data)
	return data
}
