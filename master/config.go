package master

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"cluster-balance/utils"
)

type Config struct {
	ClusterID       string
	EetcdHosts      []string
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

func (c *Config) Validate() error {
	if !utils.ValidateNodeName(c.ClusterID) {
		return errors.New("invalid cluster id")
	}
	if len(c.EetcdHosts) == 0 {
		return errors.New("invalid zookeeper config")
	}
	return nil
}