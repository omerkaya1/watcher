package config

import (
	"github.com/omerkaya1/watcher/internal/errors"
	"github.com/spf13/viper"
	"path"
)

type Config struct {
	DB    DBConf    `json:"db" yaml:"db" toml:"db"`
	Queue QueueConf `json:"queue" yaml:"queue" toml:"queue"`
}

type DBConf struct {
	Host     string `json:"host" yaml:"host" toml:"host"`
	Port     string `json:"port" yaml:"port" toml:"port"`
	Password string `json:"password" yaml:"password" toml:"password"`
	Name     string `json:"name" yaml:"name" toml:"name"`
	User     string `json:"user" yaml:"user" toml:"user"`
	SSLMode  string `json:"sslmode" yaml:"sslmode" toml:"sslmode"`
}

type QueueConf struct {
	Host      string `json:"host" yaml:"host" toml:"host"`
	Port      string `json:"port" yaml:"port" toml:"port"`
	User      string `json:"user" yaml:"user" toml:"user"`
	Password  string `json:"password" yaml:"password" toml:"password"`
	Interval  string `json:"interval" yaml:"interval" toml:"interval"`
	QueueName string `json:"queuename" yaml:"queuename" toml:"queuename"`
}

func InitConfig(cfgPath string) (*Config, error) {
	viper.SetConfigFile(cfgPath)

	cfgFileExt := path.Ext(cfgPath)
	if cfgFileExt == "" {
		return nil, errors.ErrBadConfigFile
	}
	viper.SetConfigType(cfgFileExt[1:])

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
