package Config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"short-link/pkg/logger"
)

var ErrInvalidYamlFile = errors.New("invalid yaml file")

// Config holds the app master configuration
type Config struct {
	HTTPPort    int           `yaml:"HTTP_PORT" envconfig:"HTTP_PORT"`
	GRPCPort    int           `yaml:"GRPC_PORT" envconfig:"GRPC_PORT"`
	RefererHost string        `yaml:"REFERER_HOST" envconfig:"REFERER_HOST"`
	AppMod      string        `yaml:"APP_MOD" envconfig:"APP_MOD"`
	HASHCODE    string        `yaml:"HASHCODE" envconfig:"HASHCODE"`
	DB          DB            `yaml:"DB"`
	QueueRabbit QueueRabbit   `yaml:"QueueRabbit"`
	Logger      logger.Config `yaml:"LOGGER"`
	Redis       Redis         `yaml:"REDIS"`
}

type DB struct {
	Host     string `yaml:"HOST" envconfig:"HOST"`
	Port     int    `yaml:"PORT" envconfig:"PORT"`
	User     string `yaml:"USER" `
	Password string `yaml:"PASSWORD" envconfig:"PASSWORD"`
	Dbname   string `yaml:"DBNAME" envconfig:"DBNAME"`
}

type QueueRabbit struct {
	MainQueueName string `yaml:"MAIN_QUEUE_NAME" envconfig:"MAIN_QUEUE_NAME"`
	Host          string `yaml:"HOST" envconfig:"HOST"`
	Port          int    `yaml:"PORT" envconfig:"PORT"`
	User          string `yaml:"USER" `
	Password      string `yaml:"PASSWORD" envconfig:"PASSWORD"`
}

type Redis struct {
	HOST     string `yaml:"HOST" envconfig:"HOST"`
	PORT     int    `yaml:"PORT" envconfig:"PORT"`
	USEAUTH  bool   `yaml:"USE_AUTH" envconfig:"USE_AUTH"`
	PASSWORD string `yaml:"PASSWORD" envconfig:"PASSWORD"`
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

func LoadTestConfig() (*Config, error) {
	defaultConfigFile := "../../../../../../config/config-test.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("../../config/config-%s.yaml", env)
	}
	return LoadConfig(defaultConfigFile)
}

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}
