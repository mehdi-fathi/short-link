package internal

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
)

var ErrInvalidYamlFile = errors.New("invalid yaml file")

// Config holds the app master configuration
type Config struct {
	HTTPPort    int    `yaml:"HTTP_PORT" envconfig:"HTTP_PORT"`
	GRPCPort    int    `yaml:"GRPC_PORT" envconfig:"GRPC_PORT"`
	RefererHost string `yaml:"REFERER_HOST" envconfig:"REFERER_HOST"`
	AppMod      string `yaml:"APP_MOD" envconfig:"APP_MOD"`
	HASHCODE    string `yaml:"HASHCODE" envconfig:"HASHCODE"`
	//Logger      logger.Config `yaml:"LOGGER"`
	Redis Redis `yaml:"REDIS"`
}

type Redis struct {
	HOST     []string `yaml:"HOST" envconfig:"HOST"`
	PORT     int      `yaml:"PORT" envconfig:"PORT"`
	USEAUTH  bool     `yaml:"USE_AUTH" envconfig:"USE_AUTH"`
	PASSWORD string   `yaml:"PASSWORD" envconfig:"PASSWORD"`
}

// LoadConfig loads configs form provided yaml file or overrides it with env variables
func LoadConfig(filePath string) (*Config, error) {
	cfg := Config{}
	if filePath != "" {
		err := readFile(&cfg, filePath)
		if err != nil {
			return nil, err
		}
	}
	err := readEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func readFile(cfg *Config, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return ErrInvalidYamlFile
	}
	return nil
}

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}
