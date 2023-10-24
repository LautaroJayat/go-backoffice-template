package config

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Configs struct {
	HTTP        HTTPConfig        `yaml:"HTTP"`
	DB          DBConfig          `yaml:"DB"`
	Propagation PropagationConfig `yaml:"Propagation"`
	AppName     string            `yaml:"AppName"`
}

type HTTPConfig struct {
	ReadTimeout    int    `yaml:"ReadTimeout"`
	WriteTimeout   int    `yaml:"WriteTimeout"`
	MaxHeaderBytes int    `yaml:"MaxHeaderBytes"`
	Port           string `yaml:"Port"`
}

type DBConfig struct {
	Host     string `yaml:"Host"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	Port     string `yaml:"Port"`
	SSLMode  string `yaml:"SSLMode"`
}

type PropagationConfig struct {
	Redis    Redis               `yaml:"Redis"`
	Channels PropagationChannels `yaml:"Channels"`
}

type Redis struct {
	Addr       string `yaml:"Addr"`
	User       string `yaml:"User"`
	Password   string `yaml:"Password"`
	PubTimeOut uint   `yaml:"PubTimeOut"`
}

type PropagationChannels struct {
	Customers string `yaml:"Customers"`
	Products  string `yaml:"Products"`
}

func FromYAML(f io.Reader) (*Configs, error) {
	b, err := io.ReadAll(f)

	if err != nil {
		return nil, fmt.Errorf("couldn't read config file. The error was %q", err)
	}
	cfg := &Configs{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal yaml file: %q", err)
	}
	return cfg, nil
}
