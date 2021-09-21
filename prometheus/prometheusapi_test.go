package prometheus

import (
	"fmt"
	"testing"
)

// func TestGetUpNum(t *testing.T) {
// 	nums, err := GetDeviceUpNumByJob("")
// 	if err != nil {
// 		fmt.Println("err is", err)
// 		return

// 	}
// 	fmt.Println(nums)
// }

// func TestGetTotalNum(t *testing.T) {
// 	nums, err := GetDeviceUpTotalNum("")
// 	if err != nil {
// 		fmt.Println("err is", err)
// 		return

// 	}
// 	fmt.Println(nums)
// }

// func TestGetSum(t *testing.T) {
// 	sum, err := GetMetricValuesSum("miner_available", "")
// 	fmt.Println("sum is", sum, "err is", err)
// }

// func TestGetMetricsValueDelta(t *testing.T) {
// 	delta, err := GetMetricsValueDelta("miner_block_produced", "", "3m")
// 	fmt.Println("delta is", delta, "err is", err)
// }

// func TestGetMetricValuesGroupByTimeRange(t *testing.T) {
// 	group, err := GetMetricValuesGroupByTimeRange("miner_available", "10.133.13.131", "", "1d", "1h")
// 	fmt.Println("group is", group, "err is", err)
// }

func TestGetMetricsByLocalAddr(t *testing.T) {
	resp, err := GetDeviceUpNumByJob("", "")
	if err != nil {
		fmt.Println("err is", err)
		return
	}
	fmt.Println("resp is", resp)
}
