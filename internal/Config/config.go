package Config

import (
	"flag"
	_ "flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"short-link/pkg/logger"
)

var ErrInvalidYamlFile = errors.New("invalid yaml file")

// Config holds the app master configuration
type Config struct {
	HTTPPort     int           `yaml:"HTTP_PORT" envconfig:"HTTP_PORT"`
	GRPCPort     int           `yaml:"GRPC_PORT" envconfig:"GRPC_PORT"`
	RefererHost  string        `yaml:"REFERER_HOST" envconfig:"REFERER_HOST"`
	HttpProtocol string        `yaml:"HTTP_PROTOCOL" envconfig:"HTTP_PROTOCOL"`
	AppMod       string        `yaml:"APP_MOD" envconfig:"APP_MOD"`
	HASHCODE     string        `yaml:"HASHCODE" envconfig:"HASHCODE"`
	DB           DB            `yaml:"DB"`
	QueueRabbit  QueueRabbit   `yaml:"QueueRabbit"`
	Logger       logger.Config `yaml:"LOGGER"`
	Redis        Redis         `yaml:"REDIS"`
	AppPath      string        `yaml:"APP_PATH"`
}

type DB struct {
	Driver   string `yaml:"DRIVER" envconfig:"DRIVER"`
	Host     string `yaml:"HOST" envconfig:"HOST"`
	Port     int    `yaml:"PORT" envconfig:"PORT"`
	User     string `yaml:"USER" envconfig:"USER"`
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

var ConfigHandy *Config

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

	ConfigHandy = &cfg

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

func GetBaseUrl() string {
	url := fmt.Sprintf("%s://%s", ConfigHandy.HttpProtocol, ConfigHandy.RefererHost)

	if ConfigHandy.HTTPPort > 0 {
		url = fmt.Sprintf("%s:%d", url, ConfigHandy.HTTPPort)
	}

	return url
}

func LoadConfigEnvApp() *Config {

	// Load environment variables from .env file
	env1 := []string{".env.local", ".env"}
	err := godotenv.Load(env1...)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var Config Config
	// Process environment variables into struct
	err = envconfig.Process("", &Config)
	if err != nil {
		log.Fatalf("Error processing env variables: %v", err)
	}

	ConfigHandy = &Config

	return &Config
}

func LoadConfigApp() *Config {
	// Default Config file based on the environment variable
	defaultConfigFile := "config/config-local.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("config/config-%s.yaml", env)
	}

	// Load Master Config File
	var configFile string
	flag.StringVar(&configFile, "config", defaultConfigFile, "The environment configuration file of application")
	flag.Usage = usage
	flag.Parse()

	// Loading the config file
	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Println(errors.Wrapf(err, "failed to load config: %s", "CreateService"))
	}

	if err != nil {
		log.Fatal(errors.Wrap(err, "logger error"))
	}

	return cfg
}

func usage() {
	usageStr := `
Usage: server [options]
Options:
	-c,  --config   <config file name>   Path of yaml configuration file
`
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
