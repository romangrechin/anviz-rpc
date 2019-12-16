package models

type DeviceConnectRequest struct {
	Host string `json:"host"`
}

type DeviceConnectDisconnect struct {
	Id uint32 `json:"id"`
}

type DeviceInfoResponse struct {
	Id        uint32 `json:"id"`
	FwVersion string `json:"fw_version"`
	Password  string `json:"password"`
	SleepTime uint8  `json:"sleep_time"`
	Volume    uint8  `json:"volume"`
}

type DateTime struct {
	DateTime string `json:"datetime"`
}

type StateResponse struct {
	Code uint8  `json:"code"`
	Text string `json:"text"`
}

type Records struct {
	Users        int32 `json:"users"`
	Fingerprints int32 `json:"fingerprints"`
	Passwords    int32 `json:"passwords"`
	Cards        int32 `json:"cards"`
	All          int32 `json:"all"`
	New          int32 `json:"new"`
}

type Capacity struct {
	Users        int32 `json:"users"`
	Fingerprints int32 `json:"fingerprints"`
	Records      int32 `json:"records"`
}

type StatusResponse struct {
	Records  Records       `json:"records"`
	Capacity Capacity      `json:"capacity"`
	State    StateResponse `json:"state"`
}

type Users struct {
	Id             uint64 `json:"id,omitempty"`
	Password       uint32 `json:"password"`
	CardCode       uint32 `json:"card_code"`
	Name           string `json:"name"`
	Department     uint8  `json:"department"`
	Group          uint8  `json:"group"`
	AttendanceMode uint8  `json:"attendance_mode"`
	RegisteredFp   uint16 `json:"registered_fp"`
	Keep           uint8  `json:"keep"`
	SpecialInfo    uint8  `json:"special_info"`
	IsAdmin        bool   `json:"is_admin"`
}

type RecordItem struct {
	UserId         uint64 `json:"user_id"`
	DateTime       string `json:"datetime"`
	BackupCode     uint8  `json:"backup_code"`
	Type           string `json:"type"`
	AttendanceMode uint8  `json:"attendance_mode"`
	WorkTypes      int32  `json:"work_types"`
}
