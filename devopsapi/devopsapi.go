package devopsapi

import (
	"encoding/json"
	"fmt"
	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
	etcdcli "github.com/NpoolDevOps/fbc-license-service/etcdcli"
	"github.com/NpoolRD/http-daemon"
	"golang.org/x/xerrors"
)

const devopsDomain = "devops.npool.top"

func MyDevicesByUsername(input types.MyDevicesByUsernameInput) (*types.MyDevicesOutput, error) {
	host, err := etcdcli.GetHostByDomain(devopsDomain)
	if err != nil {
		return nil, err
	}

	log.Infof(log.Fields{}, "req to http://%v%v", host, types.MyDevicesByUsernameAPI)

	resp, err := httpdaemon.R().
		SetHeader("Content-Type", "application/json").
		SetBody(input).
		Post(fmt.Sprintf("http://%v%v", host, types.MyDevicesByUsernameAPI))
	if err != nil {
		log.Errorf(log.Fields{}, "heartbeat error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, xerrors.Errorf("NON-200 return")
	}

	apiResp, err := httpdaemon.ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	output := types.MyDevicesOutput{}
	b, _ := json.Marshal(apiResp.Body)
	err = json.Unmarshal(b, &output)

	return &output, err
}
