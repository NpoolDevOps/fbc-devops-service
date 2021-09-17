package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

type MetricRangeResponse struct {
	Status string          `json:"status"`
	Data   MetricRangeData `json:"data"`
}

type MetricRangeData struct {
	Result []MetricRangeResult `json:"result"`
}

type MetricRangeResult struct {
	Metric map[string]string `json:"metric"`
	Value  [][]interface{}   `json:"values"`
}

type MetricData struct {
	Result []MetricResult `json:"result"`
}

type MetricResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type MyMetric struct {
	MetricName string            `json:"metric_name"`
	Metric     map[string]string `json:"metric"`
	Value      string            `json:"value"`
}

type Metrics struct {
	Metric []MyMetric `json:"metric"`
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

func GetMetricsByLocalAddr(localAddr string) (Metrics, error) {
	var output Metrics

	query := "{instance=\"" + localAddr + ":52379\"}"
	query = strings.Replace(url.QueryEscape(query), "+", "%20", -1)

	response, err := getQueryResponse(query, false)
	if err != nil {
		return output, err
	}

	for _, v := range response.(MetricResponse).Data.Result {
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

func GetMetricsByTime(queryTime, address, metric string) (float64, error) {

	query := "{instance=\"" + address + ":52379\"}"
	query = strings.Replace(url.QueryEscape(query), "+", "%20", -1)
	query = metric + query + "&" + "time=" + queryTime

	response, err := getQueryResponse(query, false)
	if err != nil {
		return 0, err
	}

	if len(response.(MetricResponse).Data.Result) == 0 {
		return 0, nil
	}

	result, _ := strconv.ParseFloat(response.(MetricResponse).Data.Result[0].Value[1].(string), 64)

	return result, nil

}

func GetMetricValueByAddress(address, metric string) (string, error) {
	query := "{instance=\"" + address + ":52379\"}"
	query = metric + strings.Replace(url.QueryEscape(query), "+", "%20", -1)

	response, err := getQueryResponse(query, false)
	if err != nil {
		return "", err
	}

	if len(response.(MetricResponse).Data.Result) == 0 {
		return "", nil
	}
	return response.(MetricResponse).Data.Result[0].Value[1].(string), nil

}

type DeviceMetricByTime struct {
	Diff     uint64
	Instance string
}

type AllDevicesMetricResult struct {
	MetricName string
	Results    []DeviceMetricByTime
}

func GetDevicesMetricsDiffByTime(metrics []string, startTime, endTime int64) ([]AllDevicesMetricResult, error) {
	result := []AllDevicesMetricResult{}

	for _, metric := range metrics {
		metricResult := AllDevicesMetricResult{}

		query := fmt.Sprintf("%v&start=%v&end=%v&step=%v", metric, startTime, endTime, (endTime - startTime))
		response, err := getQueryResponse(query, true)
		if err != nil {
			return nil, err
		}

		fmt.Println("response is", response)
		metricResult.MetricName = metric

		for _, data := range response.(MetricRangeResponse).Data.Result {
			dataResult := DeviceMetricByTime{}
			dataResult.Instance = strings.Split(data.Metric["instance"], ":")[0]
			valueAfter := (data.Value[1])[1].(string)
			valueAfterToUint, _ := strconv.ParseUint(valueAfter, 10, 64)
			valueBefore := (data.Value[0])[1].(string)
			valueBeforeToUint, _ := strconv.ParseUint(valueBefore, 10, 64)
			dataResult.Diff = valueAfterToUint - valueBeforeToUint

			metricResult.Results = append(metricResult.Results, dataResult)
		}

		result = append(result, metricResult)
	}

	return result, nil
}

func getQueryResponse(query string, doRange bool) (interface{}, error) {
	response := MetricResponse{}
	rangeResponse := MetricRangeResponse{}
	url := ""
	if doRange {
		url = fmt.Sprintf("http://%v/api/v1/query_range?query=%v", PrometheusSite, query)
	} else {
		url = fmt.Sprintf("http://%v/api/v1/query?query=%v", PrometheusSite, query)
	}
	fmt.Println("url is", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf(log.Fields{}, "get info from prometheus host %v error: %v", PrometheusSite, err)
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(log.Fields{}, "read resp.Body error: %v", err)
		return response, err
	}

	if doRange {
		err = json.Unmarshal(body, &rangeResponse)
		fmt.Println("err is", err, "response range is", rangeResponse)
		return rangeResponse, err
	} else {
		err = json.Unmarshal(body, &response)
		fmt.Println("err is", err, "response range is", response)
		return response, err
	}
}
