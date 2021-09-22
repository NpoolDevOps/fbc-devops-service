package fbcredis

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/EntropyPool/entropy-logger"
	types "github.com/NpoolDevOps/fbc-devops-service/types"
	etcdcli "github.com/NpoolDevOps/fbc-license-service/etcdcli"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type RedisConfig struct {
	Host string        `json:"host"`
	Ttl  time.Duration `json:"ttl"`
}

type RedisCli struct {
	config RedisConfig
	client *redis.Client
}

func NewRedisCli(config RedisConfig) *RedisCli {
	cli := &RedisCli{
		config: config,
	}

	var myConfig RedisConfig

	resp, err := etcdcli.Get(config.Host)
	if err == nil {
		err = json.Unmarshal(resp[0], &myConfig)
		if err == nil {
			cli = &RedisCli{
				config: myConfig,
			}
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr: cli.config.Host,
		DB:   0,
	})

	log.Infof(log.Fields{}, "redis ping -> %v", config.Host)
	pong, err := client.Ping().Result()
	if err != nil {
		log.Errorf(log.Fields{}, "new redis client error [%v]", err)
		return nil
	}

	if pong != "PONG" {
		log.Errorf(log.Fields{}, "redis connect failed!")
	} else {
		log.Infof(log.Fields{}, "redis connect success!")
	}

	cli.client = client

	return cli
}

var redisKeyPrefix = "fbc:devop:server:"

func (cli *RedisCli) InsertKeyInfo(keyWord string, id uuid.UUID, info interface{}, ttl time.Duration) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = cli.client.Set(fmt.Sprintf("%v:%v:%v", redisKeyPrefix, keyWord, id),
		string(b), ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cli *RedisCli) QueryDevice(cid uuid.UUID) (*types.DeviceConfig, error) {
	val, err := cli.client.Get(fmt.Sprintf("%v:device:%v", redisKeyPrefix, cid)).Result()
	if err != nil {
		return nil, err
	}
	info := &types.DeviceConfig{}
	err = json.Unmarshal([]byte(val), info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

type ClientInfo struct {
	Id          string `json:"id"`
	ClientUser  string `json:"client_user"`
	ClientSn    string `json:"client_sn"`
	Status      string `json:"status"`
	CreateTime  string `json:"create_time"`
	ModifyTime  string `json:"modify_time"`
	NetworkType string `json:"network_type"`
}

func (cli *RedisCli) QueryClient(id string) (ClientInfo, error) {
	val, err := cli.client.Get(fmt.Sprintf("fbc:license:server::client:%v", id)).Result()
	if err != nil {
		return ClientInfo{}, err
	}

	clientInfo := ClientInfo{}
	err = json.Unmarshal([]byte(val), &clientInfo)
	if err != nil {
		return ClientInfo{}, err
	}
	return clientInfo, nil
}
