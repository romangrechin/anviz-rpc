package models

import (
	"encoding/binary"
	"github.com/romangrechin/anviz-rpc/anviz/errors"
	"strconv"
)

const (
	UnicodeUsersMaxResults = 0x08
	AnsiUsersMaxResults    = 0x0C
)

type UserList struct {
	length    int
	users     []User
	IsUnicode bool
}

func (ul *UserList) Add(user User) error {
	max := AnsiUsersMaxResults
	if ul.IsUnicode {
		max = UnicodeUsersMaxResults
	}
	if len(ul.users) >= max {
		return errors.ErrTaUserListFull
	}

	ul.users = append(ul.users, user)
	ul.length = len(ul.users)
	return nil
}

func (ul *UserList) Get() []User {
	return ul.users
}

func (ul *UserList) Length() int {
	return ul.length
}

func (ul *UserList) SetIsUnicode(val bool) {
	ul.IsUnicode = val
}

func (ul *UserList) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return errors.ErrInvalidUserListData
	}

	dataLength := 30
	if ul.IsUnicode {
		dataLength = 40
	}

	ul.length = int(data[0])
	if len(data) != 1+ul.length*dataLength {
		return errors.ErrInvalidUserListData
	}

	for i := 1; i < len(data); i += dataLength {
		u := &User{}
		u.SetIsUnicode(ul.IsUnicode)
		err := u.Unmarshal(data[i : i+dataLength])
		if err != nil {
			return errors.ErrInvalidTaRecordListData
		}

		err = ul.Add(*u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ul *UserList) Marshal() []byte {
	dataLength := 30
	if ul.IsUnicode {
		dataLength = 40
	}

	buf := make([]byte, 1+ul.length*dataLength)
	buf[0] = uint8(ul.length)
	for i := 0; i < ul.length; i++ {
		copy(buf[i*40+1:], ul.users[i].Marshal())
	}
	return buf
}

type User struct {
	UserCode       uint64
	pwdLength      uint8
	Password       uint32
	CardCode       uint32
	Name           string
	Department     uint8
	Group          uint8
	AttendanceMode uint8
	RegisteredFp   uint16
	Keep           uint8
	SpecialInfo    uint8
	IsUnicode      bool
	isAdmin        bool
}

func (u *User) SetIsUnicode(val bool) {
	u.IsUnicode = val
}

func (u *User) SetIsAdmin(val bool) {
	u.isAdmin = val
}

func (u *User) IsAdmin() bool {
	return u.isAdmin
}

func (u *User) Unmarshal(data []byte) error {
	dataLength := 30
	offset := 22
	if u.IsUnicode {
		dataLength = 40
		offset = 32
	}

	if len(data) != dataLength {
		return errors.ErrInvalidUserData
	}

	if !isEmpty(data[5:8]) {
		u.Password, u.pwdLength = unpackUserPassword(data[5:8], data[offset+5])
	} else {
		u.Password, u.pwdLength = 0, 0
	}

	u.UserCode = toInt64(data[0:5])

	if !isEmpty(data[8:12]) {
		u.CardCode = binary.BigEndian.Uint32(data[8:12])
	} else {
		u.CardCode = 0
	}

	u.Name = unicodeToString(data[12:offset])
	u.Department = data[offset]
	u.Group = data[offset+1]
	u.AttendanceMode = data[offset+2]
	u.RegisteredFp = binary.BigEndian.Uint16(data[offset+3 : offset+5])
	u.Keep = data[offset+6]
	u.SpecialInfo = data[offset+7]
	u.isAdmin = u.SpecialInfo>>6 == 3
	return nil
}

func (u *User) Marshal() []byte {
	dataLength := 30
	offset := 22
	max := 10
	if u.IsUnicode {
		dataLength = 40
		offset = 32
		max = 20
	}

	data := make([]byte, dataLength)

	bufCode := fromUInt64(u.UserCode)
	copy(data[0:5], bufCode)

	bufPwd, add := []byte{0xff, 0xff, 0xff}, uint8(0x00)
	if u.Password > 0 {
		u.pwdLength = uint8(len(strconv.Itoa(int(u.Password))))
		bufPwd, add = packUserPassword(u.Password, u.pwdLength)
	}
	copy(data[5:8], bufPwd)
	data[offset+5] = add

	bufCardCode := []byte{0xff, 0xff, 0xff, 0xff}
	if u.CardCode > 0 {
		binary.BigEndian.PutUint32(bufCardCode, u.CardCode)
	}
	copy(data[8:12], bufCardCode)

	bufName := stringToUnicode(u.Name, max)
	copy(data[12:offset], bufName)

	data[offset] = u.Department
	data[offset+1] = u.Group
	data[offset+2] = u.AttendanceMode

	binary.BigEndian.PutUint16(data[offset+3:offset+5], u.RegisteredFp)
	data[offset+6] = u.Keep

	if u.isAdmin {
		u.SpecialInfo = (u.SpecialInfo & 0x3f) | (uint8(3) << 6)
	} else {
		u.SpecialInfo = (u.SpecialInfo & 0x3f) | (uint8(1) << 6)
	}

	data[offset+7] = u.SpecialInfo
	return data
}

type UserCreateResponse struct {
}

func (ucr *UserCreateResponse) Unmarshal(data []byte) error {
	return nil
}

type UserDeleteRequest struct {
	Id         uint64
	BackupCode uint8
}

func (udr *UserDeleteRequest) Marshal() []byte {
	data := make([]byte, 6)
	bufCode := fromUInt64(udr.Id)
	copy(data[0:5], bufCode)
	data[5] = udr.BackupCode
	return data
}
