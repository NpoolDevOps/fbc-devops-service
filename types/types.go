package types

import (
	"github.com/google/uuid"
)

type DeviceRegisterInput struct {
	Id          uuid.UUID `gorm:"column:id" json:"id,empty"`
	Spec        string    `gorm:"column:spec" json:"spec"`
	ParentSpec  string    `gorm:"column:parent_spec" json:"parent_spec"`
	Role        string    `gorm:"column:role" json:"role"`
	SubRole     string    `gorm:"column:sub_role" json:"sub_role"`
	Owner       string    `gorm:"column:owner" json:"owner"`
	CurrentUser string    `gorm:"column:current_user" json:"current_user"`
	Manager     string    `gorm:"column:manager" json:"manager"`
	NvmeCount   int       `gorm:"column:nvme_count" json:"nvme_count"`
	NvmeDesc    []string  `json:"nvme_desc"`
	GpuCount    int       `gorm:"column:gpu_count" json:"gpu_count"`
	GpuDesc     []string  `json:"gpu_desc"`
	MemoryCount int       `gorm:"column:memory_count" json:"memory_count"`
	MemorySize  uint64    `gorm:"column:memory_size" json:"memory_size"`
	MemoryDesc  []string  `json:"memory_desc"`
	CpuCount    int       `gorm:"column:cpu_count" json:"cpu_count"`
	CpuDesc     []string  `json:"cpu_desc"`
	HddCount    int       `gorm:"column:hdd_count" json:"hdd_count"`
	HddDesc     []string  `json:"hdd_desc"`
	LocalAddr   string    `json:"local_addr"`
	PublicAddr  string    `json:"public_addr"`
}

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
}

type DeviceReportOutput struct {
	DeviceCommonOutput
}
