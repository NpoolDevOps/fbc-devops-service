package fbcmysql

import (
	"encoding/json"
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	etcdcli "github.com/NpoolDevOps/fbc-license-service/etcdcli"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/xerrors"
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
		err = json.Unmarshal(resp, &myConfig)
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
	Id          uuid.UUID `gorm:"column:id;primary_key"`
	Spec        string    `gorm:"column:spec"`
	ParentSpec  string    `gorm:"column:parent_spec"`
	Role        string    `gorm:"column:role"`
	SubRole     string    `gorm:"column:sub_role"`
	NvmeCount   int       `gorm:"column:nvme_count"`
	NvmeDesc    string    `gorm:"column:nvme_desc"`
	GpuCount    int       `gorm:"column:gpu_count"`
	GpuDesc     string    `gorm:"column:gpu_desc"`
	MemoryCount int       `gorm:"column:memory_count"`
	MemorySize  uint64    `gorm:"column:memory_size"`
	MemoryDesc  string    `gorm:"column:memory_desc"`
	CpuCount    int       `gorm:"column:cpu_count"`
	CpuDesc     string    `gorm:"column:cpu_desc"`
	HddCount    int       `gorm:"column:hdd_count"`
	HddDesc     string    `gorm:"column:hdd_desc"`
	Maintaining bool      `gorm:"column:maintaining"`
	Offline     bool      `gorm:"column:offline"`
	Updating    bool      `gorm:"column:updating"`
	CreateTime  time.Time `gorm:"column:create_time"`
	ModifyTime  time.Time `gorm:"column:modify_time"`
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

type DeviceRuntime struct {
	Id          string    `gorm:"column:id;primary_key"`
	NvmeCount   int       `gorm:"column:nvme_count"`
	GpuCount    int       `gorm:"column:gpu_count"`
	MemoryCount int       `gorm:"column:memory_count"`
	MemorySize  uint64    `gorm:"column:memory_size"`
	ReportTime  time.Time `gorm:"column:report_time"`
}

func (cli *MysqlCli) QueryDeviceRuntime(id uuid.UUID) (*DeviceRuntime, error) {
	var info DeviceRuntime
	var count int

	cli.db.Where("id = ?", id).Find(&info).Count(&count)
	if count == 0 {
		return nil, xerrors.Errorf("cannot find any value")
	}

	return &info, nil
}

func (cli *MysqlCli) InsertElement(info interface{}) error {
	rc := cli.db.Create(&info)
	return rc.Error
}

func (cli *MysqlCli) QueryDeviceConfigs() []DeviceConfig {
	var infos []DeviceConfig
	cli.db.Find(&infos)
	return infos
}

func (cli *MysqlCli) QueryDeviceRuntimes() []DeviceRuntime {
	var infos []DeviceRuntime
	cli.db.Find(&infos)
	return infos
}
