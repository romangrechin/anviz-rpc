package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/romangrechin/anviz-rpc/anviz/device"
	"github.com/romangrechin/anviz-rpc/anviz/errors"
	devmodels "github.com/romangrechin/anviz-rpc/anviz/models"
	"github.com/romangrechin/anviz-rpc/api/models"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	ErrForbidden          = errors.New("forbidden", 403)
	ErrDeviceNotFound     = errors.New("device not found", 404)
	ErrUserNotFound       = errors.New("user not found", 405)
	ErrFpTemplateNotFound = errors.New("template not found", 406)
)

var (
	tokenAuth string
	server    *http.Server
)

type deviceListItem struct {
	Id uint32 `json:"id"`
	Ip string `json:"ip"`
}

type devices struct {
	items map[uint32]*device.Device
	hosts map[string]uint32
	mu    sync.Mutex
}

func (ds *devices) Add(host string, id uint32, dev *device.Device) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.items[id] = dev
	ds.hosts[host] = id
}

func (ds *devices) Get(id uint32) (dev *device.Device, err error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	var ok bool

	if dev, ok = ds.items[id]; ok {
		return
	}

	return nil, ErrDeviceNotFound
}

func (ds *devices) GetByHost(host string) (dev *device.Device, err error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	var (
		ok bool
		id uint32
	)

	if id, ok = ds.hosts[host]; ok {
		if dev, ok = ds.items[id]; ok {
			return
		}
	}
	delete(ds.hosts, host)

	return nil, ErrDeviceNotFound
}

func (ds *devices) RemoveByHost(host string) (err error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var (
		ok  bool
		id  uint32
		dev *device.Device
	)

	if id, ok = ds.hosts[host]; ok {
		if dev, ok = ds.items[id]; ok {
			dev.Disconnect()
		}
		delete(ds.items, id)
	} else {
		return ErrDeviceNotFound
	}
	delete(ds.hosts, host)
	return nil
}

func (ds *devices) RemoveById(id uint32) (err error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var (
		ok   bool
		dev  *device.Device
		host string
	)

	if dev, ok = ds.items[id]; ok {
		dev.Disconnect()
	} else {
		return ErrDeviceNotFound
	}
	delete(ds.items, id)

	for key, val := range ds.hosts {
		if val == id {
			host = key
			break
		}
	}

	delete(ds.hosts, host)
	return nil
}

func (ds *devices) DisconnectAll() {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if len(ds.items) == 0 {
		return
	}

	for _, dev := range ds.items {
		dev.Disconnect()
	}

	ds.items = make(map[uint32]*device.Device)
	ds.hosts = make(map[string]uint32)
}

func (ds *devices) GetAll() []deviceListItem {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var list []deviceListItem

	for _, dev := range ds.items {
		item := deviceListItem{
			Id: dev.Id(),
		}

		for ip, id := range ds.hosts {
			if id == item.Id {
				item.Ip = ip
				break
			}
		}

		list = append(list, item)
	}

	return list
}

var (
	devs *devices
)

func init() {
	devs = &devices{items: make(map[uint32]*device.Device), hosts: make(map[string]uint32)}
}

func baseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json; charset=utf-8")
		resp := &models.Response{}
		defer func() {
			json.NewEncoder(w).Encode(resp)
		}()

		key := r.Header.Get("X-API-Key")
		if key != tokenAuth {
			resp.Error = ErrForbidden
			return
		}

		ctx := context.WithValue(r.Context(), "resp", resp)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func deviceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := (r.Context().Value("resp")).(*models.Response)
		id := getDeviceId(r)

		if id == 0 {
			resp.Error = ErrDeviceNotFound
			return
		}

		dev, err := devs.Get(id)
		if err != nil {
			resp.Error = ErrDeviceNotFound
			return
		}

		ctx := context.WithValue(r.Context(), "dev", dev)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getDeviceId(r *http.Request) uint32 {
	vars := mux.Vars(r)
	if idString, ok := vars["id"]; ok {
		id, err := strconv.Atoi(idString)
		if err != nil {
			return 0
		}

		if id > math.MaxInt32 {
			return 0
		}

		return uint32(id)
	}
	return 0
}

func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := (r.Context().Value("resp")).(*models.Response)
		id := getUserId(r)

		if id == 0 {
			resp.Error = ErrUserNotFound
			return
		}

		ctx := context.WithValue(r.Context(), "user", id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getUserId(r *http.Request) uint64 {
	vars := mux.Vars(r)
	if idString, ok := vars["user_id"]; ok {
		id, err := strconv.Atoi(idString)
		if err != nil {
			return 0
		}

		if id > math.MaxInt32 {
			return 0
		}

		return uint64(id)
	}
	return 0
}

func getFpId(r *http.Request) uint8 {
	vars := mux.Vars(r)
	if idString, ok := vars["fp_id"]; ok {
		id, err := strconv.Atoi(idString)
		if err != nil {
			return 0
		}

		if id > 2 {
			return 0
		}

		return uint8(id)
	}
	return 0
}

func connect(w http.ResponseWriter, r *http.Request) {
	resp := (r.Context().Value("resp")).(*models.Response)

	connectReq := &models.DeviceConnectRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(connectReq)
	if err != nil {
		log.Println(err)
		resp.Error = errors.New(err.Error(), 10000)
		return
	}

	dev, err := devs.GetByHost(connectReq.Host)
	if err == nil {
		_, err = dev.GetDateTime()
		if err == nil || err == errors.ErrDeviceIsBusy {
			resp.Data = &models.DeviceConnectDisconnect{Id: dev.Id()}
			return
		}
		_ = devs.RemoveByHost(connectReq.Host)
	}

	dev = device.New(connectReq.Host, 10*time.Second, 2*time.Second)
	err = dev.Connect()
	if err != nil {
		resp.Error = err
		return
	}

	devs.Add(connectReq.Host, dev.Id(), dev)
	resp.Data = &models.DeviceConnectDisconnect{Id: dev.Id(), Code: dev.Code(), BiometricType: dev.BiometricType()}
}

func disconnect(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)
	resp.Error = devs.RemoveById(dev.Id())
}

func getDateTime(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	dt, err := dev.GetDateTime()
	if err != nil {
		resp.Error = err
		return
	}

	resp.Data = &models.DateTime{DateTime: dt.Format("02-01-2006 15:04:05")}
}

func setDateTime(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	req := &models.DateTime{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(req)
	if err != nil {
		log.Println(err)
		resp.Error = errors.New(err.Error(), 10000)
		return
	}

	datetime, err := time.Parse("02-01-2006 15:04:05", req.DateTime)
	if err != nil {
		log.Println(err)
		resp.Error = errors.New(err.Error(), 10000)
		return
	}

	err = dev.SetDateTime(devmodels.DateTime{Time: datetime})
	if err != nil {
		resp.Error = err
	}
}

func state(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	state := &models.StateResponse{
		Code: dev.State(),
	}

	switch state.Code {
	case 0:
		state.Text = "disconnected"
	case 1:
		state.Text = "connected"
	case 2:
		state.Text = "busy"
	default:
		state.Text = "unknown"
	}
	resp.Data = state
}

func status(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	state := &models.StatusResponse{}

	ri, err := dev.GetRecordInfo()
	if err != nil {
		resp.Error = err
		return
	}

	state.Records.Users = ri.Users
	state.Records.New = ri.New
	state.Records.All = ri.All
	state.Records.Cards = ri.Cards
	state.Records.Fingerprints = ri.FingerPrints
	state.Records.Passwords = ri.Passwords

	cap, err := dev.GetCapacity()
	if err != nil {
		resp.Error = err
		return
	}

	state.Capacity.Records = cap.Records
	state.Capacity.Fingerprints = cap.Fingerprints
	state.Capacity.Users = cap.Users

	state.State.Code = dev.State()
	switch state.State.Code {
	case 0:
		state.State.Text = "disconnected"
	case 1:
		state.State.Text = "connected"
	case 2:
		state.State.Text = "busy"
	default:
		state.State.Text = "unknown"
	}
	resp.Data = state
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	var users []models.Users

	first := true
	count := 0
	max := dev.MaxUserResult()

	for {
		ul, err := dev.GetUsers(first, false)
		if err != nil {
			resp.Error = err
			return
		}

		count = ul.Length()
		for _, u := range ul.Get() {
			user := &models.Users{
				Id:             u.UserCode,
				Password:       u.Password,
				CardCode:       u.CardCode,
				Name:           u.Name,
				Department:     u.Department,
				Group:          u.Group,
				AttendanceMode: u.AttendanceMode,
				RegisteredFp:   u.RegisteredFp,
				Keep:           u.Keep,
				SpecialInfo:    u.SpecialInfo,
				IsAdmin:        u.IsAdmin(),
			}

			users = append(users, *user)
			resp.Data = users
		}

		first = false

		if count < max {
			break
		}
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	req := &models.Users{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(req)
	if err != nil {
		log.Println(err)
		resp.Error = errors.New(err.Error(), 10000)
		return
	}

	// TODO: Validate

	list := devmodels.UserList{}
	list.SetIsUnicode(true)
	list.SetIsC3(dev.IsC3())
	user := &devmodels.User{
		UserCode:       req.Id,
		Password:       req.Password,
		CardCode:       req.CardCode,
		Name:           req.Name,
		Department:     req.Department,
		Group:          req.Group,
		AttendanceMode: req.AttendanceMode,
		RegisteredFp:   req.RegisteredFp,
		Keep:           req.Keep,
		SpecialInfo:    req.SpecialInfo,
		IsUnicode:      true,
		IsC3:           dev.IsC3(),
	}
	user.SetIsAdmin(req.IsAdmin)
	err = list.Add(*user)
	if err != nil {
		resp.Error = err
		return
	}

	err = dev.SetUsers(list)
	if err != nil {
		resp.Error = err
	}
}

func modifyUser(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)
	userId := (r.Context().Value("user")).(uint64)

	req := &models.Users{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(req)
	if err != nil {
		log.Println(err)
		resp.Error = errors.New(err.Error(), 10000)
		return
	}

	// TODO: Validate

	list := devmodels.UserList{}
	list.SetIsUnicode(true)
	list.SetIsC3(dev.IsC3())
	user := &devmodels.User{
		UserCode:       userId,
		Password:       req.Password,
		CardCode:       req.CardCode,
		Name:           req.Name,
		Department:     req.Department,
		Group:          req.Group,
		AttendanceMode: req.AttendanceMode,
		RegisteredFp:   req.RegisteredFp,
		Keep:           req.Keep,
		SpecialInfo:    req.SpecialInfo,
		IsUnicode:      true,
		IsC3:           dev.IsC3(),
	}
	user.SetIsAdmin(req.IsAdmin)
	err = list.Add(*user)
	if err != nil {
		resp.Error = err
		return
	}

	err = dev.SetUsers(list)
	if err != nil {
		resp.Error = err
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)
	userId := (r.Context().Value("user")).(uint64)

	var backup uint8 = 0xff

	_ = r.ParseForm()
	backupString := r.Form.Get("backup")
	if backupString != "" {
		b, err := strconv.Atoi(backupString)
		if err != nil {
			resp.Error = errors.New("query: invalid field \"backup\": "+err.Error(), 4001)
			return
		}
		backup = uint8(b)
	}

	req := devmodels.UserDeleteRequest{
		Id:         userId,
		BackupCode: backup,
	}

	err := dev.DeleteUser(req)
	if err != nil {
		resp.Error = err
	}
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	newOnly := false

	_ = r.ParseForm()
	if r.Form.Get("new_only") == "1" {
		newOnly = true
	}

	var records []models.RecordItem

	recs, err := dev.GetTaRecords(newOnly)
	if err != nil {
		resp.Error = err
		return
	}

	for _, rec := range recs {
		item := models.RecordItem{
			UserId:         rec.UserCode,
			DateTime:       rec.DateTime.Format("02-01-2006 15:04:05"),
			BackupCode:     rec.BackupCode,
			AttendanceMode: rec.AttendanceMode(),
			WorkTypes:      rec.WorkTypes,
		}

		if rec.IsOpen() {
			item.Type = "in"
		} else {
			item.Type = "out"
		}

		records = append(records, item)
	}

	resp.Data = records
}

func clearNewRecords(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)

	var clearType uint8 = 1
	count := 0

	_ = r.ParseForm()
	countString := r.Form.Get("count")
	if countString != "" {
		count, _ = strconv.Atoi(countString)
	}

	if count > 0 {
		clearType = 2
	}

	total, err := dev.ClearRecord(clearType, int32(count))
	if err != nil {
		resp.Error = err
		return
	}

	resp.Data = total
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	resp := (r.Context().Value("resp")).(*models.Response)
	resp.Data = devs.GetAll()
}

func getFpTemplate(w http.ResponseWriter, r *http.Request) {
	dev := (r.Context().Value("dev")).(*device.Device)
	resp := (r.Context().Value("resp")).(*models.Response)
	userId := (r.Context().Value("user")).(uint64)
	fpId := getFpId(r)

	if fpId == 0 {
		resp.Error = ErrFpTemplateNotFound
		return
	}

	data, err := dev.DownloadFpTemplate(userId, fpId)
	if err != nil {
		resp.Error = ErrFpTemplateNotFound
		return
	}

	if len(data) == 0 {
		resp.Error = ErrFpTemplateNotFound
		return
	}

	resp.Data = base64.StdEncoding.EncodeToString(data)
}

func index(w http.ResponseWriter, r *http.Request) {
	var data = `<html>
<head>
<title>anviz-rpc</title>
</head>
<body>
anviz-rpc
</body>
</html>
`

	w.Write([]byte(data))
}

func RunServer(address string, token string) error {
	if server != nil {
		return errors.New("server already running: " + server.Addr)
	}
	tokenAuth = token
	r := mux.NewRouter().Methods("GET", "POST").Subrouter()
	r.Use(baseMiddleware)
	r.HandleFunc("/", index)
	r.HandleFunc("/connect", connect)
	r.HandleFunc("/devices", getDevices)

	s := r.PathPrefix("/{id:[0-9]+}").Methods("GET", "POST").Subrouter()
	s.Use(deviceMiddleware)
	s.HandleFunc("/disconnect", disconnect).Methods("GET")
	s.HandleFunc("/datetime", getDateTime).Methods("GET")
	s.HandleFunc("/datetime", setDateTime).Methods("POST")
	s.HandleFunc("/state", state).Methods("GET")
	s.HandleFunc("/status", status).Methods("GET")
	s.HandleFunc("/records", getRecords).Methods("GET")
	s.HandleFunc("/clear_new_records", clearNewRecords).Methods("GET")
	u := s.PathPrefix("/users").Methods("GET", "POST").Subrouter()
	u.HandleFunc("/", getUsers).Methods("GET")
	u.HandleFunc("/add", addUser).Methods("POST")
	ur := u.PathPrefix("/{user_id:[0-9]+}").Methods("GET", "POST").Subrouter()
	ur.Use(userMiddleware)
	ur.HandleFunc("/modify", modifyUser).Methods("POST")
	ur.HandleFunc("/delete", deleteUser).Methods("GET")
	ur.HandleFunc("/fp/{fp_id:[0-9]+}", getFpTemplate).Methods("GET")

	server = &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stoppped with error: %s\n", err.Error())
		}
		server = nil
	}()
	time.Sleep(500 * time.Millisecond)
	return nil
}

func StopServer() error {
	if server == nil {
		return errors.New("server is not running")
	}
	defer func() {
		server = nil
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	return server.Shutdown(ctx)
}
