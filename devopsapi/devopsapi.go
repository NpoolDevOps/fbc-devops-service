package devopsapi

import (
	"encoding/json"
	"fmt"
	"time"
	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
	etcdcli "github.com/NpoolDevOps/fbc-license-service/etcdcli"
	"github.com/NpoolRD/http-daemon"
	"golang.org/x/xerrors"
)

const devopsDomain = "devops.npool.top"

func MyDevicesByUsername(input types.MyDevicesByUsernameInput, useDomain bool) (*types.MyDevicesOutput, error) {
	var host string
	var err error
	var scheme string

	if useDomain {
		scheme = "https"
		host = devopsDomain
	} else {
		host, err = etcdcli.GetHostByDomain(devopsDomain)
		if err != nil {
			return nil, err
		}
		scheme = "http"
	}

	log.Infof(log.Fields{}, "req to %v://%v%v", scheme, host, types.MyDevicesByUsernameAPI)

	resp, err := httpdaemon.Cli().SetTimeout(30*time.Minute).R().
		SetHeader("Content-Type", "application/json").
		SetBody(input).
		Post(fmt.Sprintf("%v://%v%v", scheme, host, types.MyDevicesByUsernameAPI))
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
