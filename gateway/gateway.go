package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
)

const (
	PrometheusSite = "106.14.125.55:9988"
)

type MetricResponse struct {
	Status string     `json:"status"`
	Data   MetricData `json:"data"`
}

type MetricData struct {
	Result []MetricResult `json:"result"`
}

type MetricResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
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
		resp, err := http.Get(fmt.Sprintf("http://%v/api/v1/query?query=%v", PrometheusSite, metric))
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
			// mymetric.Instance = strings.TrimSpace(strings.Split(v.Metric.Instance, ":")[0])
			// mymetric.Job = v.Metric.Job
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
		resp, err := http.Get(fmt.Sprintf("http://%v/api/v1/query_range?query=%v&start=%v&end=%v&step=%v", PrometheusSite, metric, startTime, endTime, step))
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

func GetMetricsByLocalAddr(localAddr string) (Metrics, error) {
	var output Metrics
	result := MetricResponse{}
	query := "{instance=\"" + localAddr + ":52379\"}"
	query = strings.Replace(url.QueryEscape(query), "+", "%20", -1)
	resp, err := http.Get(fmt.Sprintf("http://%v/api/v1/query?query=%v", PrometheusSite, query))
	if err != nil {
		log.Errorf(log.Fields{}, "get info from prometheus error: %v", err)
		return Metrics{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(log.Fields{}, "read resp body error: %v", err)
		return Metrics{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Errorf(log.Fields{}, "fail to unmarshal: %v", err)
		return Metrics{}, err
	}

	for _, v := range result.Data.Result {
		out := MyMetric{}
		if strings.HasPrefix(v.Metric["__name__"], "go_") || strings.HasPrefix(v.Metric["__name__"], "process_") || strings.HasPrefix(v.Metric["__name__"], "promhttp_") || v.Metric["__name__"] == "miner_seal_sector_task_progress" {
			continue
		}
		out.Value = v.Value[1].(string)
		out.MetricName = v.Metric["__name__"]
		metric := make(map[string]string)
		for k, vv := range v.Metric {
			metric[k] = vv
		}
		out.Metric = metric
		output.Metric = append(output.Metric, out)
	}
	return output, nil
}

type MyMetric struct {
	MetricName string            `json:"metric_name"`
	Metric     map[string]string `json:"metric"`
	Value      string            `json:"value"`
}

type Metrics struct {
	Metric []MyMetric `json:"metric"`
}

func GetMetricsByTime(queryTime, address, metric string) (float64, error) {
	response := MetricResponse{}
	query := "{instance=\"" + address + ":52379\"}"
	query = strings.Replace(url.QueryEscape(query), "+", "%20", -1)
	query = metric + query + "&" + "time=" + queryTime

	resp, err := http.Get(fmt.Sprintf("http://%v/api/v1/query?query=%v", PrometheusSite, query))
	if err != nil {
		log.Errorf(log.Fields{}, "get info from prometheus error: %v", err)
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	if len(response.Data.Result) == 0 {
		return 0, nil
	}

	result, _ := strconv.ParseFloat(response.Data.Result[0].Value[1].(string), 64)

	return result, nil

}
