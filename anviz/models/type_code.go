package models

type TypeCode struct {
	Code string
}

func (tc *TypeCode) Unmarshal(data []byte) error {
	var buf []byte
	for i := range data {
		if data[i] == 0x00 {
			break
		}

		buf = append(buf, data[i])
	}
	tc.Code = string(buf)
	return nil
}
