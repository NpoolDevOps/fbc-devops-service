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
	NvmeDesc    []string  `gorm:"column:nvme_desc" json:"nvme_desc"`
	GpuCount    int       `gorm:"column:gpu_count" json:"gpu_count"`
	GpuDesc     []string  `gorm:"column:gpu_desc" json:"gpu_desc"`
	MemoryCount int       `gorm:"column:memory_count" json:"memory_count"`
	MemorySize  uint64    `gorm:"column:memory_size" json:"memory_size"`
	MemoryDesc  []string  `gorm:"column:memory_desc" json:"memory_desc"`
	CpuCount    int       `gorm:"column:cpu_count" json:"cpu_count"`
	CpuDesc     []string  `gorm:"column:cpu_desc" json:"cpu_desc"`
	HddCount    int       `gorm:"column:hdd_count" json:"hdd_count"`
	HddDesc     []string  `gorm:"column:hdd_desc" json:"hdd_desc"`
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
