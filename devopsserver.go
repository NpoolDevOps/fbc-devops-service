package main

import (
	"encoding/json"
	log "github.com/EntropyPool/entropy-logger"
	devopsmysql "github.com/NpoolDevOps/fbc-devops-service/mysql"
	devopsredis "github.com/NpoolDevOps/fbc-devops-service/redis"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
	"github.com/NpoolRD/http-daemon"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

	log.Infof(log.Fields{}, "start http daemon at %v", s.config.Port)
	httpdaemon.Run(s.config.Port)
	return nil
}

func (s *DevopsServer) DeviceRegisterRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	b, err := ioutil.ReadAll(req.Body)

	input := types.DeviceRegisterInput{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err.Error(), -1
	}

	config := devopsmysql.DeviceConfig{}
	config.Id = input.Id
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
	config.CreateTime = time.Now()
	config.ModifyTime = time.Now()

	err = s.mysqlClient.InsertDeviceConfig(config)
	if err != nil {
		return nil, err.Error(), -2
	}

	return nil, "", 0
}

func (s *DevopsServer) DeviceReportRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	return nil, "", 0
}

func (s *DevopsServer) DeviceMaintainRequest(w http.ResponseWriter, req *http.Request) (interface{}, string, int) {
	return nil, "", 0
}