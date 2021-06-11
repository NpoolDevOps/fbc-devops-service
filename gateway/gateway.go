package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
)

type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Result []Result `json:"result"`
}

type Result struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}

type Metric struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

// type Outresp struct {
// 	metricName string
// 	metric     []MyMetric
// }

// type MyMetric struct {
// 	instance string
// 	job      string
// 	value    string
// }

func GetMetrics(metrics []string) ([]types.Outresp, error) {
	var output []types.Outresp
	for _, metric := range metrics {
		result := Response{}
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
			mymetric.Instance = v.Metric.Instance
			mymetric.Job = v.Metric.Job
			mymetric.Value = v.Value[1].(string)

			outputResp.Metric = append(outputResp.Metric, mymetric)
			outputResp.MetricName = metric
		}

		output = append(output, outputResp)
	}
	return output, nil
}
