package models

import "github.com/romangrechin/anviz-rpc/anviz/errors"

type RecordInfo struct {
	Users        int32
	FingerPrints int32
	Passwords    int32
	Cards        int32
	All          int32
	New          int32
}

func (ri *RecordInfo) Unmarshal(data []byte) error {
	if len(data) != 18 {
		return errors.ErrInvalidRecordInfoData
	}

	ri.Users = toInt(data[0:3])
	ri.FingerPrints = toInt(data[3:6])
	ri.Passwords = toInt(data[6:9])
	ri.Cards = toInt(data[9:12])
	ri.All = toInt(data[12:15])
	ri.New = toInt(data[15:18])

	return nil
}
