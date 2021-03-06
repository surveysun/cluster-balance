package comm

import (
	"cluster-balance/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Config struct {
	ClusterID string   `json:"cluster_id"`
	EtcdHosts []string `json:"etcd_hosts"`
}

const BaseNode = "/cluster_balance"

func LoadConfig(fn string) (*Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) GetClusterPath() string {
	return fmt.Sprintf("%s/%s", BaseNode, c.ClusterID)
}

func (c *Config) GetNodesPath() string {
	return fmt.Sprintf("%s/%s/nodes", BaseNode, c.ClusterID)
}

func (c *Config) GetResourcesPath() string {
	return fmt.Sprintf("%s/%s/resources", BaseNode, c.ClusterID)
}

func (c *Config) GetStatusPath() string {
	return fmt.Sprintf("%s/%s/status", BaseNode, c.ClusterID)
}

func (c *Config) GetMasterPath() string {
	return fmt.Sprint("%s/%s/master", BaseNode, c.ClusterID)
}

func (c *Config) Validate() error {
	if !utils.ValidateNodeName(c.ClusterID) {
		return errors.New("invalid cluster id")
	}
	if len(c.EtcdHosts) == 0 {
		return errors.New("invalid etcd config")
	}
	return nil
}
