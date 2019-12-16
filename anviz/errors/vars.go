package errors

var (
	ErrOK               = New("success", 0)
	ErrCouldNotConnect  = New("could not connect to device", 1)
	ErrAlreadyConnected = New("already connected to device", 2)
	ErrConnectionClosed = New("connection closed", 3)
	ErrConnectionWrite  = New("connection write error", 4)
	ErrConnectionRead   = New("connection read error", 5)
	ErrDeviceIsBusy     = New("device is busy", 6)
)

var (
	ErrInvalidResponseData     = New("invalid response data", 100)
	ErrInvalidResponseAck      = New("invalid response ack", 101)
	ErrInvalidResponseChecksum = New("invalid response checksum", 102)
	ErrDeviceOperationFailed   = New("operation failed", 103)
	ErrDeviceOperationCanceled = New("operation canceled", 104)
)

var (
	ErrInvalidDeviceInfoData   = New("invalid device info data", 1000)
	ErrInvalidDateTimeData     = New("invalid datetime data", 1001)
	ErrInvalidRecordInfoData   = New("invalid record info data", 1002)
	ErrInvalidTaRecordListData = New("invalid T&A record list data", 1003)
	ErrInvalidTaRecordData     = New("invalid T&A record data", 1004)
	ErrTaRecordListFull        = New("T&A record list is full", 1005)
	ErrInvalidTimestampData    = New("invalid timestamp data", 1006)
	ErrInvalidUserListData     = New("invalid user list data", 1007)
	ErrInvalidUserData         = New("invalid user data", 1008)
	ErrTaUserListFull          = New("user list is full", 1009)
	ErrInvalidCapacityData     = New("invalid capacity data", 1010)
)
