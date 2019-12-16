package models

import (
	"encoding/binary"
	"github.com/romangrechin/anviz-rpc/anviz/errors"
	"time"
)

var (
	time2000 time.Time
)

type DateTime struct {
	time.Time
}

type Timestamp struct {
	time.Time
}

func (ts *Timestamp) Unmarshal(data []byte) error {
	timestamp := binary.BigEndian.Uint32(data)
	ts.Time = time2000.Add(time.Duration(timestamp) * time.Second)
	return nil
}

func (dt *DateTime) Unmarshal(data []byte) error {
	if len(data) != 6 {
		return errors.ErrInvalidDateTimeData
	}

	t := time.Time{}
	dt.Time = t.AddDate(2000+int(data[0])-1, int(data[1])-1, int(data[2])-1).
		Add(time.Duration(data[3]) * time.Hour).
		Add(time.Duration(data[4]) * time.Minute).
		Add(time.Duration(data[5]) * time.Second)
	return nil
}

func (dt *DateTime) Marshal() []byte {
	buf := make([]byte, 6)
	buf[0] = uint8(dt.Year() - 2000)
	buf[1] = uint8(dt.Month())
	buf[2] = uint8(dt.Day())
	buf[3] = uint8(dt.Hour())
	buf[4] = uint8(dt.Minute())
	buf[5] = uint8(dt.Second())
	return buf
}

func init() {
	time2000 = time.Time{}
	_ = time2000.UnmarshalText([]byte("2000-01-01T00:00:00Z"))
}
