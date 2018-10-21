package config

import (
	"github.com/kindle_server/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type GlobalConfig struct {
	MySQL      *MySQLConfig   `yaml:"mysql"`
	LifeTime   types.Duration `yaml:"lifetime"`
	GcInterval types.Duration `yaml:"gcInterval"`
	Split      string         `yaml:"split"`
}

func LoadFile(filename string) (*GlobalConfig, []byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	cfg := &GlobalConfig{}
	err = yaml.UnmarshalStrict([]byte(content), cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, content, nil

}

func (c *GlobalConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain GlobalConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}
