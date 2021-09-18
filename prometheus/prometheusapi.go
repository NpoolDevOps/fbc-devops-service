package prometheus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/EntropyPool/entropy-logger"
	"github.com/NpoolDevOps/fbc-devops-service/types"
	"golang.org/x/xerrors"
)

const (
	PrometheusSite = "106.14.125.55:9988"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Result []Result `json:"result"`
}

type Result struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
	Values [][]interface{}   `json:"values"`
}

func getQueryResponse(query string) (Response, error) {
	response := Response{}
	url := fmt.Sprintf("http://%v/api/v1/query?query=%v", PrometheusSite, query)
	resp, err := http.Get(url)
	if err != nil {
		return Response{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}
	if len(response.Data.Result) == 0 {
		return Response{}, xerrors.Errorf("there is no result")
	}
	return response, nil
}

func queryByUser(customerName string) string {
	if customerName != "" {
		return fmt.Sprintf("user=\"%v\",", customerName)
	}
	return ""
}

func queryByNetworkType(networkType string) string {
	if networkType != "" {
		return fmt.Sprintf("networktype=\"%v\",", networkType)
	}
	return ""
}

type UpNum struct {
	Job        string `json:"job"`
	UpDevice   uint64 `json:"up_device"`
	DownDevice uint64 `json:"down_device"`
}

func GetDeviceUpNumByJob(customerName, networkType string) ([]UpNum, error) {
	jobs := []string{types.FullNode, types.MinerNode, types.FullMinerNode, types.WorkerNode, types.StorageNode}

	var output []UpNum
	for _, job := range jobs {
		query := fmt.Sprintf("count_values(\"count\",up{instance=~\".*:52379\",job=\"%v\",%v%v})", job, queryByUser(customerName), queryByNetworkType(networkType))

		response, err := getQueryResponse(query)
		upNum := UpNum{}
		upNum.Job = job
		upNum.DownDevice = 0
		upNum.UpDevice = 0
		if err != nil {
			output = append(output, upNum)
			continue
		}

		for _, result := range response.Data.Result {
			num, _ := strconv.ParseUint(result.Value[1].(string), 10, 64)
			if result.Metric["count"] == "1" {
				upNum.UpDevice = num
			} else if result.Metric["count"] == "0" {
				upNum.DownDevice = num
			}
		}
		output = append(output, upNum)
	}
	return output, nil
}

func GetDeviceUpTotalNum(customerName, networkType string) (uint64, error) {
	query := fmt.Sprintf("count(up{instance=~\".*:52379\",%v%v})", queryByUser(customerName), queryByNetworkType(networkType))

	response, err := getQueryResponse(query)
	if err != nil {
		return 0, err
	}

	num, _ := strconv.ParseUint(response.Data.Result[0].Value[1].(string), 10, 64)
	return num, nil
}

func GetMetricValuesSum(metricName, customerName, networkType string) (float64, error) {
	query := fmt.Sprintf("sum(%v{%v%v})", metricName, queryByUser(customerName), queryByNetworkType(networkType))

	response, err := getQueryResponse(query)
	if err != nil {
		return 0, err
	}
	sum, _ := strconv.ParseFloat(response.Data.Result[0].Value[1].(string), 64)
	return sum, nil

}

func GetMetricsValueDelta(metricName, customerName, timeRange, networkType string) (float64, error) {
	query := fmt.Sprintf("sum(delta(%v{%v%v}[%v]))", metricName, queryByUser(customerName), queryByNetworkType(networkType), timeRange)

	response, err := getQueryResponse(query)
	if err != nil {
		return 0, err
	}
	delta, _ := strconv.ParseFloat(response.Data.Result[0].Value[1].(string), 64)
	return delta, nil
}

type ValueGroup struct {
	MetricName string                   `json:"metric_name"`
	Instance   string                   `json:"instance"`
	Values     []map[string]interface{} `json:"values"`
}

func GetMetricValuesGroupByTimeRange(metricName, instance, customerName, timeRange, timeStep string) (ValueGroup, error) {
	query := fmt.Sprintf("%v{instance=\"%v:52379\",%v}[%v:%v]", metricName, instance, queryByUser(customerName), timeRange, timeStep)
	timeLayout := "2006-01-02 15:04:05"
	response, err := getQueryResponse(query)
	if err != nil {
		return ValueGroup{}, err
	}
	valueGroup := ValueGroup{}
	valueGroup.Instance = instance
	valueGroup.MetricName = metricName
	for _, value := range response.Data.Result[0].Values {
		timeStamp := time.Unix(int64(value[0].(float64)), 0).Format(timeLayout)
		val, _ := strconv.ParseFloat(value[1].(string), 64)
		myValue := make(map[string]interface{})
		myValue["time"] = timeStamp
		myValue["value"] = val
		valueGroup.Values = append(valueGroup.Values, myValue)
	}

	return valueGroup, nil
}

type MyMetric struct {
	MetricName string            `json:"metric_name"`
	Metric     map[string]string `json:"metric"`
	Value      string            `json:"value"`
}

type Metrics struct {
	Metric []MyMetric `json:"metric"`
}

func GetMetricsByLocalAddr(localAddr string) (Metrics, error) {
	query := fmt.Sprintf("{instance=\"%v:52379\"}", localAddr)
	response, err := getQueryResponse(query)
	if err != nil {
		return Metrics{}, err
	}

	var metrics Metrics
	for _, result := range response.Data.Result {
		if strings.HasPrefix(result.Metric["__name__"], "go_") || strings.HasPrefix(result.Metric["__name__"], "process_") || strings.HasPrefix(result.Metric["__name__"], "promhttp_") || result.Metric["__name__"] == "miner_seal_sector_task_progress" {
			continue
		}

		metric := MyMetric{}
		metric.Value = result.Value[1].(string)
		metric.MetricName = result.Metric["__name__"]
		metr := make(map[string]string)
		for k, v := range result.Metric {
			metr[k] = v
		}
		metric.Metric = metr
		metrics.Metric = append(metrics.Metric, metric)
	}

	return metrics, nil
}

type MetricValueDelta struct {
	MetricName string
	Delta      float64
}

func GetMetricsValueDeltaByAddress(metrics []string, localAddr, customerName, timeRange string) ([]MetricValueDelta, error) {
	output := []MetricValueDelta{}
	for _, metric := range metrics {
		var metricValueDelta MetricValueDelta
		query := fmt.Sprintf("delta(%v{instance=\"%v:52379\",%v}[%v])", metric, localAddr, queryByUser(customerName), timeRange)
		response, err := getQueryResponse(query)
		metricValueDelta.MetricName = metric
		if err != nil {
			metricValueDelta.Delta = 0
			log.Errorf(log.Fields{}, "fail to get %v delta, error is %v", metric, err)
			continue
		}
		metricValueDelta.Delta, _ = strconv.ParseFloat(response.Data.Result[0].Value[1].(string), 64)
		output = append(output, metricValueDelta)
	}

	return output, nil
}
