package models

import (
	"github.com/romangrechin/anviz-rpc/anviz/errors"
)

const (
	TaRecordsMaxResults = 0x19
)

type TaRecordList struct {
	length  int
	records []TaRecord
}

func (trl *TaRecordList) Add(record TaRecord) error {
	if len(trl.records) >= TaRecordsMaxResults {
		return errors.ErrTaRecordListFull
	}

	trl.records = append(trl.records, record)
	return nil
}

func (trl *TaRecordList) Get() []TaRecord {
	return trl.records
}

func (trl *TaRecordList) Length() int {
	return trl.length
}

func (trl *TaRecordList) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return errors.ErrInvalidTaRecordListData
	}

	trl.length = int(data[0])
	if len(data) != 1+trl.length*14 {
		return errors.ErrInvalidTaRecordListData
	}

	for i := 1; i < len(data); i += 14 {
		tr := &TaRecord{}
		err := tr.Unmarshal(data[i : i+14])
		if err != nil {
			return errors.ErrInvalidTaRecordListData
		}

		err = trl.Add(*tr)
		if err != nil {
			return err
		}
	}

	return nil
}

type TaRecord struct {
	UserCode   uint64
	DateTime   Timestamp
	BackupCode uint8
	RecordType uint8
	WorkTypes  int32
}

func (tr *TaRecord) IsOpen() bool {
	return (tr.RecordType & 1) == 1
}

func (tr *TaRecord) AttendanceMode() uint8 {
	return tr.RecordType & 0xf0 >> 4
}

func (tr *TaRecord) Unmarshal(data []byte) error {
	if len(data) != 14 {
		return errors.ErrInvalidTaRecordData
	}

	tr.UserCode = toInt64(data[0:5])
	tr.DateTime = Timestamp{}
	_ = tr.DateTime.Unmarshal(data[5:9])
	tr.BackupCode = data[9]
	tr.RecordType = data[10]
	tr.WorkTypes = toInt(data[11:14])
	return nil
}
