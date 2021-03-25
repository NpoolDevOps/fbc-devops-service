package devopsmysql

import (
	"encoding/json"
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	etcdcli "github.com/NpoolDevOps/fbc-license-service/etcdcli"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/xerrors"
	"strings"
	"time"
)

type MysqlConfig struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	DbName string `json:"db"`
}

type MysqlCli struct {
	config MysqlConfig
	url    string
	db     *gorm.DB
}

func NewMysqlCli(config MysqlConfig) *MysqlCli {
	cli := &MysqlCli{
		config: config,
		url: fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
			config.User, config.Passwd, config.Host, config.DbName),
	}

	var myConfig MysqlConfig

	resp, err := etcdcli.Get(config.Host)
	if err == nil {
		err = json.Unmarshal(resp[0], &myConfig)
		if err == nil {
			myConfig.DbName = config.DbName
			cli = &MysqlCli{
				config: myConfig,
				url: fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
					myConfig.User, myConfig.Passwd, myConfig.Host, myConfig.DbName),
			}
		}
	}

	log.Infof(log.Fields{}, "open mysql db %v", cli.url)
	db, err := gorm.Open("mysql", cli.url)
	if err != nil {
		log.Errorf(log.Fields{}, "cannot open %v: %v", cli.url, err)
		return nil
	}

	log.Infof(log.Fields{}, "successful to create mysql db %v", cli.url)
	db.SingularTable(true)
	cli.db = db

	return cli
}

func (cli *MysqlCli) Delete() {
	cli.db.Close()
}

type DeviceConfig struct {
	Maintaining bool      `gorm:"column:maintaining"`
	Offline     bool      `gorm:"column:offline"`
	Updating    bool      `gorm:"column:updating"`
	CreateTime  time.Time `gorm:"column:create_time"`
	ModifyTime  time.Time `gorm:"column:modify_time"`
	NvmeDesc    string    `gorm:"column:nvme_desc"`
	GpuDesc     string    `gorm:"column:gpu_desc"`
	MemoryDesc  string    `gorm:"column:memory_desc"`
	CpuDesc     string    `gorm:"column:cpu_desc"`
	HddDesc     string    `gorm:"column:hdd_desc"`
	Id          uuid.UUID `gorm:"column:id"`
	Spec        string    `gorm:"column:spec"`
	ParentSpec  string    `gorm:"column:parent_spec"`
	Role        string    `gorm:"column:role"`
	SubRole     string    `gorm:"column:sub_role"`
	Owner       string    `gorm:"column:owner"`
	CurrentUser string    `gorm:"column:current_user"`
	Manager     string    `gorm:"column:manager"`
	NvmeCount   int       `gorm:"column:nvme_count"`
	GpuCount    int       `gorm:"column:gpu_count"`
	MemoryCount int       `gorm:"column:memory_count"`
	MemorySize  uint64    `gorm:"column:memory_size"`
	CpuCount    int       `gorm:"column:cpu_count"`
	HddCount    int       `gorm:"column:hdd_count"`
	OsSpec      string    `gorm:"column:os_spec"`
}

func (cli *MysqlCli) QueryDeviceConfig(id uuid.UUID) (*DeviceConfig, error) {
	var info DeviceConfig
	var count int

	cli.db.Where("id = ?", id).Find(&info).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any value")
	}

	return &info, nil
}

func (cli *MysqlCli) InsertDeviceConfig(info DeviceConfig) error {
	couldBeUpdated := false

	oldInfo, err := cli.QueryDeviceConfig(info.Id)
	if err == nil && oldInfo != nil {
		s := strings.Split(oldInfo.ParentSpec, ",")
		if oldInfo.ParentSpec == "" {
			s = []string{info.ParentSpec}
			oldInfo.ParentSpec = strings.Join(s, ",")
			couldBeUpdated = true
		} else if !strings.Contains(oldInfo.ParentSpec, info.ParentSpec) {
			s = append(s, info.ParentSpec)
			oldInfo.ParentSpec = strings.Join(s, ",")
			couldBeUpdated = true
		}
	}

	var updateInfo *DeviceConfig

	if couldBeUpdated {
		updateInfo = oldInfo
	}

	if oldInfo == nil {
		updateInfo = &info
	} else {
		if oldInfo.Maintaining {
			info.ParentSpec = oldInfo.ParentSpec
			updateInfo = &info
		}
	}

	if updateInfo == nil {
		return xerrors.Errorf("invalid operation without maintaining mode")
	}

	if couldBeUpdated {
		return cli.db.Save(updateInfo).Error
	}

	return cli.db.Create(updateInfo).Error
}

func (cli *MysqlCli) QueryDeviceConfigs() ([]DeviceConfig, error) {
	var infos []DeviceConfig
	rc := cli.db.Find(&infos)
	if rc.Error != nil {
		return nil, rc.Error
	}
	return infos, nil
}

func (cli *MysqlCli) QueryDeviceConfigsByUser(username string) ([]DeviceConfig, error) {
	var infos []DeviceConfig
	rc := cli.db.Where("username = ?", username).Find(&infos)
	if rc.Error != nil {
		return nil, rc.Error
	}
	return infos, nil
}
