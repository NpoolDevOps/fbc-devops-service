package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/EntropyPool/entropy-logger"
	authapi "github.com/NpoolDevOps/fbc-auth-service/authapi"
	authtypes "github.com/NpoolDevOps/fbc-auth-service/types"
	"github.com/NpoolDevOps/fbc-devops-service/gateway"
	devopsmysql "github.com/NpoolDevOps/fbc-devops-service/mysql"
	devopsredis "github.com/NpoolDevOps/fbc-devops-service/redis"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
	licapi "github.com/NpoolDevOps/fbc-license-service/licenseapi"
	lictypes "github.com/NpoolDevOps/fbc-license-service/types"
	httpdaemon "github.com/NpoolRD/http-daemon"
	"github.com/google/uuid"
)

type DevopsConfig struct {
	RedisCfg devopsredis.RedisConfig `json:"redis"`
	MysqlCfg devopsmysql.MysqlConfig `json:"mysql"`
	Port     int                     `json:"port"`
}

type DevopsServer struct {
	config      DevopsConfig
	authText    string
	redisClient *devopsredis.RedisCli
	mysqlClient *devopsmysql.MysqlCli
}

func NewDevopsServer(configFile string) *DevopsServer {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot read file %v: %v", configFile, err)
		return nil
	}

	config := DevopsConfig{}
	err = json.Unmarshal(buf, &config)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot parse file %v: %v", configFile, err)
		return nil
	}

	log.Infof(log.Fields{}, "create redis cli: %v", config.RedisCfg)
	redisCli := devopsredis.NewRedisCli(config.RedisCfg)
	if redisCli == nil {
		log.Errorf(log.Fields{}, "cannot create redis client %v: %v", config.RedisCfg, err)
		return nil
	}

	log.Infof(log.Fields{}, "create mysql cli: %v", config.MysqlCfg)
	mysqlCli := devopsmysql.NewMysqlCli(config.MysqlCfg)
	if mysqlCli == nil {
		log.Errorf(log.Fields{}, "cannot create mysql client %v: %v", config.MysqlCfg, err)
		return nil
	}

	server := &DevopsServer{
		config:      config,
		authText:    types.DevopsAuthText,
		redisClient: redisCli,
		mysqlClient: mysqlCli,
	}

	log.Infof(log.Fields{}, "successful to create devops server")

	return server
}

func (s *DevopsServer) Run() error {
	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.DeviceRegisterAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.DeviceRegisterRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.DeviceReportAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.DeviceReportRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.DeviceMaintainAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.DeviceMaintainRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.MyDevicesAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.MyDevicesByAuthRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.MyDevicesByAuthAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.MyDevicesByAuthRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.MyDevicesByUsernameAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.MyDevicesByUsernameRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.DevopsAlertMgrAddressAPI,
		Method:   "GET",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.DevopsAlertMgrAddressGetRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.DevopsAlertMgrAddressAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.DevopsAlertMgrAddressPostRequest(w, req)
		},
	})

	httpdaemon.RegisterRouter(httpdaemon.HttpRouter{
		Location: types.MyDevicesMetricsAPI,
		Method:   "POST",
		Handler: func(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
			return s.DevicesMetricsRequest(w, req)
		},
	})

	log.Infof(log.Fields{}, "start http daemon at %v", s.config.Port)
	httpdaemon.Run(s.config.Port)
	return nil
}

func (s *DevopsServer) DeviceRegisterRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.DeviceRegisterInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	clientInfo, err := licapi.ClientInfoBySpec(lictypes.ClientInfoBySpecInput{
		Spec: input.Spec,
	})
	if err != nil {
		return nil, err.Error(), -3
	}

	if input.Role == "" {
		return nil, "role is must", -4
	}

	valid, err := s.mysqlClient.ValidateRole(input.Role)
	if err != nil {
		return nil, err.Error(), -5
	}

	if !valid {
		return nil, "role is not valid", -6
	}

	config := devopsmysql.DeviceConfig{}
	config.Id = clientInfo.Id
	config.Spec = input.Spec
	config.ParentSpec = input.ParentSpec
	config.Role = input.Role
	config.SubRole = input.SubRole
	config.Owner = input.Owner
	config.CurrentUser = input.CurrentUser
	config.Manager = input.Manager
	config.NvmeCount = input.NvmeCount
	config.NvmeDesc = strings.Join(input.NvmeDesc, ",")
	config.GpuCount = input.GpuCount
	config.GpuDesc = strings.Join(input.GpuDesc, ",")
	config.MemoryCount = input.MemoryCount
	config.MemorySize = input.MemorySize
	config.MemoryDesc = strings.Join(input.MemoryDesc, ",")
	config.CpuCount = input.CpuCount
	config.CpuDesc = strings.Join(input.CpuDesc, ",")
	config.HddCount = input.HddCount
	config.HddDesc = strings.Join(input.HddDesc, ",")
	config.EthernetCount = input.EthernetCount
	config.EthernetDesc = strings.Join(input.EthernetDesc, ",")
	config.CreateTime = time.Now()
	config.ModifyTime = time.Now()

	input.Id = clientInfo.Id
	err = s.redisClient.InsertKeyInfo("device", input.Id, input, 2*time.Hour)
	if err != nil {
		return nil, err.Error(), -7
	}

	err = s.mysqlClient.InsertDeviceConfig(config)
	if err != nil {
		return nil, err.Error(), -8
	}

	output := types.DeviceRegisterOutput{}
	output.Id = clientInfo.Id

	return output, "", 0
}

func (s *DevopsServer) DeviceReportRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.DeviceReportInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	device, err := s.redisClient.QueryDevice(input.Id)
	if err != nil {
		return nil, err.Error(), -3
	}

	device.NvmeCount = input.NvmeCount
	device.GpuCount = input.GpuCount
	device.MemoryCount = input.MemoryCount
	device.MemorySize = input.MemorySize
	device.HddCount = input.HddCount
	device.LocalAddr = input.LocalAddr
	device.PublicAddr = input.PublicAddr

	err = s.redisClient.InsertKeyInfo("device", input.Id, device, 2*time.Hour)
	if err != nil {
		return nil, err.Error(), -3
	}

	return nil, "", 0
}

func (s *DevopsServer) DeviceMaintainRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.MaintainingInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.AuthCode == "" {
		return nil, "auth code is must", -3
	}

	user, err := authapi.UserInfo(authtypes.UserInfoInput{
		AuthCode: input.AuthCode,
	})
	if err != nil {
		return nil, err.Error(), -5
	}

	if !user.SuperUser {
		return nil, "permission denied", -6
	}

	err = s.mysqlClient.SetDeviceMaintaining(input.DeviceID, input.Maintaining)
	if err != nil {
		return nil, err.Error(), -7
	}

	return nil, "", 0
}

func (s *DevopsServer) MyDevicesByAuthRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.MyDevicesByAuthInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.AuthCode == "" {
		return nil, "auth code is must", -3
	}

	user, err := authapi.UserInfo(authtypes.UserInfoInput{
		AuthCode: input.AuthCode,
	})
	if err != nil {
		return nil, err.Error(), -4
	}

	return s.myDevicesByUserInfo(user)
}

func (s *DevopsServer) MyDevicesByUsernameRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}

	input := types.MyDevicesByUsernameInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.Username == "" {
		return nil, "username is must", -3
	}

	if input.Password == "" {
		return nil, "password is must", -4
	}

	output, err := authapi.Login(authtypes.UserLoginInput{
		Username: input.Username,
		Password: input.Password,
		AppId:    uuid.MustParse("00000002-0002-0002-0002-000000000002"),
	})
	if err != nil {
		return nil, err.Error(), -5
	}

	user, err := authapi.UserInfo(authtypes.UserInfoInput{
		AuthCode: output.AuthCode,
	})
	if err != nil {
		return nil, err.Error(), -6
	}

	return s.myDevicesByUserInfo(user)
}

func (s *DevopsServer) myDevicesByUserInfo(user *authtypes.UserInfoOutput) (interface{}, string, int) {
	infos := []devopsmysql.DeviceConfig{}
	var err error

	if user.SuperUser {
		infos, err = s.mysqlClient.QueryDeviceConfigs()
	} else {
		infos, err = s.mysqlClient.QueryDeviceConfigsByUser(user.Username)
	}

	if err != nil {
		return nil, err.Error(), -6
	}

	output := types.MyDevicesOutput{}
	for _, info := range infos {
		oInfo := types.DeviceAttribute{}

		oInfo.Id = info.Id
		oInfo.Spec = info.Spec
		oInfo.ParentSpec = strings.Split(info.ParentSpec, ",")
		oInfo.Role = info.Role
		oInfo.SubRole = info.SubRole
		oInfo.Owner = info.Owner
		oInfo.CurrentUser = info.CurrentUser
		oInfo.Manager = info.Manager
		oInfo.NvmeCount = info.NvmeCount
		oInfo.NvmeDesc = strings.Split(info.NvmeDesc, ",")
		oInfo.GpuCount = info.GpuCount
		oInfo.GpuDesc = strings.Split(info.GpuDesc, ",")
		oInfo.MemoryCount = info.MemoryCount
		oInfo.MemorySize = info.MemorySize
		oInfo.MemoryDesc = strings.Split(info.MemoryDesc, ",")
		oInfo.CpuCount = info.CpuCount
		oInfo.CpuDesc = strings.Split(info.CpuDesc, ",")
		oInfo.HddCount = info.HddCount
		oInfo.HddDesc = strings.Split(info.HddDesc, ",")
		oInfo.OsSpec = info.OsSpec
		oInfo.Maintaining = info.Maintaining
		oInfo.Offline = info.Offline

		device, err := s.redisClient.QueryDevice(info.Id)
		if err != nil {
			if !user.SuperUser {
				return nil, err.Error(), -7
			}
		} else {
			oInfo.RuntimeNvmeCount = device.NvmeCount
			oInfo.RuntimeGpuCount = device.GpuCount
			oInfo.RuntimeMemoryCount = device.MemoryCount
			oInfo.RuntimeMemorySize = device.MemorySize
			oInfo.RuntimeHddCount = device.HddCount
			oInfo.LocalAddr = device.LocalAddr
			oInfo.PublicAddr = device.PublicAddr
		}

		output.Devices = append(output.Devices, oInfo)
	}

	return output, "", 0
}

func (s *DevopsServer) DevopsAlertMgrAddressPostRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	return nil, "NOT IMPLEMENTED NOW", -1
}

func (s *DevopsServer) DevopsAlertMgrAddressGetRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	return nil, "NOT IMPLEMENTED NOW", -1
}

func (s *DevopsServer) DevicesMetricsRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err.Error(), -1
	}
	input := types.MetricInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -2
	}

	if input.AuthCode == "" {
		return nil, "auth code is must", -3
	}

	_, err = authapi.UserInfo(authtypes.UserInfoInput{
		AuthCode: input.AuthCode,
	})
	if err != nil {
		return nil, err.Error(), -4
	}

	output, err := gateway.GetMetrics(input.Metrics)
	if err != nil {
		return nil, err.Error(), -5
	}

	return types.MetricOutput{
		MetricsValue: output,
	}, "", 0
}
