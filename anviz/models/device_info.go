package models

import (
	"github.com/romangrechin/anviz-rpc/anviz/errors"
)

type DeviceInfo struct {
	FwVersion    string
	pwdLength    uint8
	Password     uint32
	SleepTime    uint8
	Volume       uint8
	Language     uint8
	DateFormat   uint8
	TimeFormat   uint8
	AttState     uint8
	LanguageFlag uint8
	CmdVersion   uint8
}

func (di *DeviceInfo) Unmarshal(data []byte) error {
	if len(data) != 18 {
		return errors.ErrInvalidDeviceInfoData
	}
	di.FwVersion = string(data[0:8])
	di.Password, di.pwdLength = unpackDevicePassword(data[8:11])
	di.SleepTime = data[11]
	di.Volume = data[12]
	di.Language = data[13]
	di.DateFormat = data[14] >> 4
	di.TimeFormat = data[14] & 0xf
	di.AttState = data[15]
	di.LanguageFlag = data[16]
	di.CmdVersion = data[17]
	return nil
}
