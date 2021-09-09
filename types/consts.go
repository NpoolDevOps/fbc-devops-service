package types

const (
	DeviceRegisterAPI              = "/api/v0/device/register"
	DeviceReportAPI                = "/api/v0/device/report"
	DeviceMaintainAPI              = "/api/v0/device/maintain"
	MyDevicesByAuthAPI             = "/api/v0/device/minebyauth"
	MyDevicesAPI                   = "/api/v0/device/mine"
	MyDevicesByUsernameAPI         = "/api/v0/device/minebyuser"
	DevopsAlertMgrAddressAPI       = "/api/v0/device/alertmgraddr"
	DevopsAuthText                 = "FBC DevOps Server - @Copyright NPool COP."
	MyDevicesMetricsAPI            = "/api/v0/device/metrics"
	DeviceMetricsDataAPI           = "/api/v0/device/metricsdata"
	MinerDeviceListAPI             = "/api/v0/miner/device/list"
	DeviceMetricsByAddressAPI      = "/api/v0/device/metrics/by/address"
	DeviceMetricByTimeAPI          = "/api/v0/device/metric/by/time"
	DeviceMetricValueDiffByTimeAPI = "/api/v0/device/metric/value/diff/by/time"
	GetAllDevicesNumAPI            = "/api/v0/get/devices/num"
	GetDeviceBlockInfosAPI         = "/api/v0/get/device/block/infos"
)

const (
	FullNode      = "fullnode"
	MinerNode     = "miner"
	FullMinerNode = "fullminer"
	WorkerNode    = "worker"
	StorageNode   = "storage"
)
