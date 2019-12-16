package models

type ClearRequest struct {
	ClearType uint8
	Count     int32
}

func (cr *ClearRequest) Marshal() []byte {
	buf := make([]byte, 4)
	if cr.ClearType < 0 || cr.ClearType > 2 {
		cr.ClearType = 1
	}

	if cr.ClearType == 2 {
		bufCount := fromInt(cr.Count)
		copy(buf[1:], bufCount)
	}

	buf[0] = cr.ClearType
	return buf
}

type ClearResponse struct {
	Count int32
}

func (cr *ClearResponse) Unmarshal(data []byte) error {
	cr.Count = toInt(data)
	return nil
}
