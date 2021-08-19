package types

import (
	"github.com/google/uuid"
)

type DeviceRegisterInput struct {
	Id            uuid.UUID `gorm:"column:id" json:"id,empty"`
	Spec          string    `gorm:"column:spec" json:"spec"`
	ParentSpec    string    `gorm:"column:parent_spec" json:"parent_spec"`
	Role          string    `gorm:"column:role" json:"role"`
	SubRole       string    `gorm:"column:sub_role" json:"sub_role"`
	Owner         string    `gorm:"column:owner" json:"owner"`
	CurrentUser   string    `gorm:"column:current_user" json:"current_user"`
	Manager       string    `gorm:"column:manager" json:"manager"`
	NvmeCount     int       `gorm:"column:nvme_count" json:"nvme_count"`
	NvmeDesc      []string  `json:"nvme_desc"`
	GpuCount      int       `gorm:"column:gpu_count" json:"gpu_count"`
	GpuDesc       []string  `json:"gpu_desc"`
	MemoryCount   int       `gorm:"column:memory_count" json:"memory_count"`
	MemorySize    uint64    `gorm:"column:memory_size" json:"memory_size"`
	MemoryDesc    []string  `json:"memory_desc"`
	CpuCount      int       `gorm:"column:cpu_count" json:"cpu_count"`
	CpuDesc       []string  `json:"cpu_desc"`
	HddCount      int       `gorm:"column:hdd_count" json:"hdd_count"`
	HddDesc       []string  `json:"hdd_desc"`
	EthernetCount int       `gorm:"column:ethernet_count" json:"ethernet_count"`
	EthernetDesc  []string  `gorm:"column:ethernet_desc" json:"ethernet_desc"`
	OsSpec        string    `gorm:"column:os_spec" json:"os_spec"`
	LocalAddr     string    `json:"local_addr"`
	PublicAddr    string    `json:"public_addr"`
	Versions      []string  `json:"versions"`
}

type DeviceConfig = DeviceRegisterInput

type DeviceCommonOutput struct {
	Id uuid.UUID `json:"id"`
}

type DeviceRegisterOutput struct {
	DeviceCommonOutput
}

type DeviceReportInput struct {
	Id          uuid.UUID `json:"id"`
	NvmeCount   int       `json:"nvme_count"`
	GpuCount    int       `json:"gpu_count"`
	MemoryCount int       `json:"memory_count"`
	MemorySize  uint64    `json:"memory_size"`
	HddCount    int       `gorm:"column:hdd_count" json:"hdd_count"`
	LocalAddr   string    `json:"local_addr"`
	PublicAddr  string    `json:"public_addr"`
}

type DeviceReportOutput struct {
	DeviceCommonOutput
}

type MyDevicesByAuthInput struct {
	AuthCode string `json:"auth_code"`
}

type MyDevicesByUsernameInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeviceAttribute struct {
	DeviceConfig
	RuntimeNvmeCount   int      `json:"runtime_nvme_count"`
	RuntimeGpuCount    int      `json:"runtime_gpu_count"`
	RuntimeMemoryCount int      `json:"runtime_memory_count"`
	RuntimeMemorySize  uint64   `json:"runtime_memory_size"`
	RuntimeHddCount    int      `json:"runtime_hdd_count"`
	ParentSpec         []string `json:"parent_spec"`
	Maintaining        bool     `json:"maintaining"`
	Offline            bool     `json:"offline"`
}

type MyDevicesOutput struct {
	Devices []DeviceAttribute `json:"devices"`
}

type MaintainingInput struct {
	AuthCode    string    `json:"auth_code"`
	Maintaining bool      `json:"maintaining"`
	DeviceID    uuid.UUID `json:"device_id"`
}

type MetricInput struct {
	Metrics  []string `json:"metrics"`
	AuthCode string   `json:"auth_code"`
}

type Outresp struct {
	MetricName string     `json:"metric_name"`
	Metric     []MyMetric `json:"metric"`
}

type MyMetric struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
	Value    string `json:"value"`
}

type MetricOutput struct {
	MetricsValue []Outresp `json:"metrics_value"`
}

type DeviceMetricsDataInput struct {
	AuthCode  string   `json:"auth_code"`
	Metrics   []string `json:"metrics"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Step      string   `json:"step"`
}

type InstanceData struct {
	Instance string        `json:"instance"`
	Job      string        `json:"job"`
	Value    []interface{} `json:"value"`
}

type MetricData struct {
	MetricName    string         `json:"metric_name"`
	InstanceDatas []InstanceData `json:"instance_datas"`
}

type DeviceMetricsDataOutput struct {
	MetricDatas []MetricData  `json:"metric_datas"`
	Date        []interface{} `json:"date"`
}

type MinerDeviceListInput struct {
	AuthCode string `json:"auth_code"`
}

type MinerDeviceListOutput = MyDevicesOutput

type MetricsByAddr struct {
	AuthCode string `json:"auth_code"`
	Address  string `json:"address"`
}

type MetricByTimeInput struct {
	AuthCode  string   `json:"auth_code"`
	QueryTime string   `json:"query_time"`
	Metric    string   `json:"metric"`
	Addresses []string `json:"addresses"`
}

type MetricByTimeOutput struct {
	Values []Value `json:"values"`
}

type Value struct {
	Address string  `json:"address"`
	Value   float64 `json:"value"`
}

type DeviceMetricValueDiffByTimeInput struct {
	AuthCode  string   `json:"auth_code"`
	Metric    string   `json:"metric"`
	Addresses []string `json:"addresses"`
	BeginTime string   `json:"begin_time"`
	EndTime   string   `json:"end_time"`
}

type GetAllDevicesNumInput struct {
}

type GetAllDevicesNumOutput struct {
	All  AllDeviceNum  `json:"all"`
	Down DownDeviceNum `json:"down"`
	Up   UpDeviceNum   `json:"up"`
}

type AllDeviceNum struct {
	WorkerNumber    uint64 `json:"worker_number"`
	MinerNumber     uint64 `json:"miner_number"`
	FullminerNumber uint64 `json:"fullminer_number"`
	FullnodeNumber  uint64 `json:"fullnode_number"`
	StorageNumber   uint64 `json:"storage_number"`
}

type UpDeviceNum struct {
	WorkerUpNumber    uint64 `json:"worker_up_number"`
	MinerUpNumber     uint64 `json:"miner_up_number"`
	FullminerUpNumber uint64 `json:"fullminer_up_number"`
	FullnodeUpNumber  uint64 `json:"fullnode_up_number"`
	StorageUpNumber   uint64 `json:"storage_up_number"`
}

type DownDeviceNum struct {
	WorkerDownNumber    uint64 `json:"worker_down_number"`
	MinerDownNumber     uint64 `json:"miner_down_number"`
	FullminerDownNumber uint64 `json:"fullminer_down_number"`
	FullnodeDownNumber  uint64 `json:"fullnode_down_number"`
	StorageDownNumber   uint64 `json:"storage_down_number"`
}
