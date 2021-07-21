package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
)

type MetricResponse struct {
	Status string     `json:"status"`
	Data   MetricData `json:"data"`
}

type MetricData struct {
	Result []MetricResult `json:"result"`
}

type MetricResult struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}
type MetricDataResponse struct {
	Status string         `json:"status"`
	Data   MetricDataData `json:"data"`
}

type MetricDataData struct {
	Result []MetricDataResult `json:"result"`
}

type MetricDataResult struct {
	Metric Metric          `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type Metric struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

func GetMetrics(metrics []string) ([]types.Outresp, error) {
	var output []types.Outresp
	for _, metric := range metrics {
		result := MetricResponse{}
		resp, err := http.Get(fmt.Sprintf("http://47.99.107.242:9090/api/v1/query?query=%v", metric))
		if err != nil {
			log.Errorf(log.Fields{}, "get info from prometheus err: %v", err)
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		outputResp := types.Outresp{}
		dataresult := result.Data.Result
		for _, v := range dataresult {
			mymetric := types.MyMetric{}
			mymetric.Instance = strings.TrimSpace(strings.Split(v.Metric.Instance, ":")[0])
			mymetric.Job = v.Metric.Job
			mymetric.Value = v.Value[1].(string)

			outputResp.Metric = append(outputResp.Metric, mymetric)
			outputResp.MetricName = metric
		}

		output = append(output, outputResp)
	}
	return output, nil
}

func GetMetricsData(metrics []string, startTime, endTime, step string) (types.DeviceMetricsDataOutput, error) {
	var output []types.MetricData
	var dateArr []interface{}
	timeTemplate := "2006-01-02-15:04:05"

	for _, metric := range metrics {
		result := MetricDataResponse{}
		resp, err := http.Get(fmt.Sprintf("http://47.99.107.242:9090/api/v1/query_range?query=%v&start=%v&end=%v&step=%v", metric, startTime, endTime, step))
		if err != nil {
			log.Errorf(log.Fields{}, "get info from prometheus err: %v", err)
			return types.DeviceMetricsDataOutput{}, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return types.DeviceMetricsDataOutput{}, err
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return types.DeviceMetricsDataOutput{}, err
		}

		metricData := types.MetricData{}
		for _, v := range result.Data.Result {
			var date []interface{}
			instanceData := types.InstanceData{}
			instanceData.Instance = strings.TrimSpace(strings.Split(v.Metric.Instance, ":")[0])
			instanceData.Job = v.Metric.Job
			for _, vv := range v.Values {
				date = append(date, time.Unix(int64(vv[0].(float64)), 0).Format(timeTemplate))
				instanceData.Value = append(instanceData.Value, vv[1])
			}
			dateArr = date
			metricData.InstanceDatas = append(metricData.InstanceDatas, instanceData)
		}
		metricData.MetricName = metric

		output = append(output, metricData)
	}
	return types.DeviceMetricsDataOutput{
		MetricDatas: output,
		Date:        dateArr,
	}, nil
}
