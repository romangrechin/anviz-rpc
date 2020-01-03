package device

import (
	"github.com/romangrechin/anviz-rpc/anviz/errors"
	"github.com/romangrechin/anviz-rpc/anviz/models"
	"time"
)

const (
	cmdDeviceInfo         = 0x30
	cmdGetDateTime        = 0x38
	cmdSetDateTime        = 0x39
	cmdGetRecordInfo      = 0x3C
	cmdGetTaRecords       = 0x40
	cmdGetC3Users         = 0x22
	cmdSetC3Users         = 0x23
	cmdGetUsers           = 0x72
	cmdSetUsers           = 0x73
	cmdDelUser            = 0x4C
	cmdGetFactoryInfoCode = 0x4A
	cmdGetCapacity        = 0x5D
	cmdClearRecord        = 0x4E
	cmdGetDeviceTypeCode  = 0x48
)

type commandResponse interface {
	Unmarshal(data []byte) error
}

type Device struct {
	conn             *connection
	address          string
	deviceInfo       *models.DeviceInfo
	isUnicode        bool
	connTimeout      time.Duration
	readWriteTimeout time.Duration
	isBusy           bool
	abort            chan struct{}
	typeCode         string
}

func (d *Device) Id() uint32 {
	if d.conn != nil {
		return d.conn.id
	}

	return 0
}

func (d *Device) IsUnicode() bool {
	return d.isUnicode
}

func (d *Device) IsC3() bool {
	return d.typeCode == "C3" || d.typeCode == "C2"
}

func (d *Device) IsBusy() bool {
	return d.isBusy
}

func (d *Device) SetUnicode(val bool) {
	d.isUnicode = val
}

func (d *Device) State() uint8 {
	if d.isBusy {
		return BUSY
	}
	return d.conn.status
}

func (d *Device) Connect() (err error) {
	// TODO: detect version (ANSI or UNICODE)
	d.SetUnicode(true)
	d.conn, err = newConnection(d.address)
	if err != nil {
		return err
	}

	_, err = d.GetInfo()
	if err != nil {
		d.Disconnect()
		return err
	}

	d.GetTypeCode()

	d.abort = make(chan struct{})

	return nil
}

func (d *Device) Disconnect() {
	if d.conn != nil {
		d.conn.Close()
		d.conn = nil
		close(d.abort)
	}
}

func (d *Device) GetTypeCode() (*models.TypeCode, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	m := &models.TypeCode{}
	err := d.request(cmdGetDeviceTypeCode, m, nil)
	if err != nil {
		d.deviceInfo = nil
		return nil, err
	}
	d.typeCode = m.Code
	return m, nil
}

func (d *Device) GetInfo() (*models.DeviceInfo, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	m := &models.DeviceInfo{}
	err := d.request(cmdDeviceInfo, m, nil)
	if err != nil {
		d.deviceInfo = nil
		return nil, err
	}
	d.deviceInfo = m
	return m, nil
}

func (d *Device) GetDateTime() (*models.DateTime, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	m := &models.DateTime{}
	err := d.request(cmdGetDateTime, m, nil)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Device) SetDateTime(m models.DateTime) error {
	if d.isBusy {
		return errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	err := d.request(cmdSetDateTime, nil, m.Marshal())
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) GetRecordInfo() (*models.RecordInfo, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	m := &models.RecordInfo{}
	err := d.request(cmdGetRecordInfo, m, nil)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Device) GetTaRecords(newRec bool) ([]models.TaRecord, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	var result []models.TaRecord

	firstPacket := true
	count := 0

	for {
		select {
		case <-d.abort:
			return nil, errors.ErrDeviceOperationCanceled
		default:
		}
		recList, err := d.getTaRecords(firstPacket, newRec, false)
		if err != nil {
			return nil, err
		}

		count = recList.Length()
		for _, rec := range recList.Get() {
			result = append(result, rec)
		}

		firstPacket = false

		if count < models.TaRecordsMaxResults {
			break
		}
	}

	return result, nil
}

func (d *Device) getTaRecords(firstPacket, newOnly, resendLastPacket bool) (*models.TaRecordList, error) {
	data := []byte{0x0, models.TaRecordsMaxResults}
	if resendLastPacket {
		data[0] = 0x10
	} else if firstPacket {
		if newOnly {
			data[0] = 0x2
		} else {
			data[0] = 0x1
		}
	}
	m := &models.TaRecordList{}
	err := d.request(cmdGetTaRecords, m, data)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Device) MaxRecordsResult() int {
	return models.TaRecordsMaxResults
}

func (d *Device) ClearRecord(clearType uint8, count int32) (int32, error) {
	// clearType=0 - clear all records
	// clearType=1 - clear all new records
	// clearType=2 - clear new records, where records count = count

	creq := &models.ClearRequest{
		ClearType: clearType,
		Count:     count,
	}
	cresp := &models.ClearResponse{}
	err := d.request(cmdClearRecord, cresp, creq.Marshal())
	if err != nil {
		return 0, err
	}
	return cresp.Count, nil
}

func (d *Device) GetAllUsers() ([]models.User, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	var result []models.User

	firstPacket := true
	count := 0

	var max uint8 = models.AnsiUsersMaxResults
	if d.isUnicode {
		max = models.UnicodeUsersMaxResults
	}

	for {
		select {
		case <-d.abort:
			return nil, errors.ErrDeviceOperationCanceled
		default:
		}
		uList, err := d.GetUsers(firstPacket, false)
		if err != nil {
			return nil, err
		}

		count = uList.Length()
		for _, rec := range uList.Get() {
			result = append(result, rec)
		}

		firstPacket = false

		if count < int(max) {
			break
		}
	}

	return result, nil
}

func (d *Device) GetUsers(firstPacket, resendLastPacket bool) (*models.UserList, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	var max uint8 = models.AnsiUsersMaxResults
	if d.isUnicode {
		max = models.UnicodeUsersMaxResults
	}
	data := []byte{0x0, max}
	if resendLastPacket {
		data[0] = 0x10
	} else if firstPacket {
		data[0] = 0x1
	}
	m := &models.UserList{}
	m.SetIsUnicode(d.isUnicode)
	var cmd uint8 = cmdGetUsers
	if d.IsC3() {
		m.SetIsC3(true)
		cmd = cmdGetC3Users
	}
	err := d.request(cmd, m, data)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Device) SetUsers(m models.UserList) error {
	if d.isBusy {
		return errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	resp := &models.UserCreateResponse{}

	var cmd uint8 = cmdSetUsers
	if d.IsC3() {
		cmd = cmdSetC3Users
	}

	err := d.request(cmd, resp, m.Marshal())
	return err
}

func (d *Device) DeleteUser(m models.UserDeleteRequest) error {
	if d.isBusy {
		return errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	err := d.request(cmdDelUser, nil, m.Marshal())
	return err
}

func (d *Device) MaxUserResult() int {
	var max uint8 = models.AnsiUsersMaxResults
	if d.isUnicode {
		max = models.UnicodeUsersMaxResults
	}
	return int(max)
}

func (d *Device) GetFactoryInfoCode() (*models.FactoryInfo, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	m := &models.FactoryInfo{}
	err := d.request(cmdGetFactoryInfoCode, m, nil)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Device) GetCapacity() (*models.Capacity, error) {
	if d.isBusy {
		return nil, errors.ErrDeviceIsBusy
	}
	d.isBusy = true
	defer func() { d.isBusy = false }()

	m := &models.Capacity{}
	err := d.request(cmdGetCapacity, m, nil)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Device) request(cmd uint8, obj commandResponse, data []byte) error {
	if d.conn == nil {
		return errors.ErrConnectionClosed
	}

	resp, err := d.conn.send(cmd, data)
	if err != nil {
		return err
	}

	if obj == nil {
		return nil
	}

	return obj.Unmarshal(resp)
}

func New(address string, connTimeout, readWriteTimeout time.Duration) *Device {
	return &Device{address: address, connTimeout: connTimeout, readWriteTimeout: readWriteTimeout}
}
