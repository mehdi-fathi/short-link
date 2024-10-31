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
	HTTPPort     int    `envconfig:"HTTP_PORT"`
	HTTPPortUrl  int    `envconfig:"HTTP_PORT_URL"`
	GRAFANAPort  int    `envconfig:"GRAFANA_PORT"`
	GRPCPort     int    `envconfig:"GRPC_PORT"`
	RefererHost  string `envconfig:"REFERER_HOST"`
	HttpProtocol string `envconfig:"HTTP_PROTOCOL"`
	AppMod       string `envconfig:"APP_MOD"`
	HASHCODE     string `envconfig:"HASHCODE"`
	DB           DB
	QueueRabbit  QueueRabbit
	Graylog      logger.Graylog `envconfig:"LOGGER_GRAYLOG"`
	Redis        REDIS          `envconfig:"REDIS"`
	AppPath      string         `envconfig:"APP_PATH"`
}

type DB struct {
	Driver   string `envconfig:"DRIVER"`
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	Dbname   string `envconfig:"DBNAME"`
}

type QueueRabbit struct {
	MainQueueName string `envconfig:"MAIN_QUEUE_NAME"`
	Host          string `envconfig:"HOST"`
	Port          int    `envconfig:"PORT"`
	User          string
	Password      string `envconfig:"PASSWORD"`
}

type REDIS struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	UseAuth  bool   `envconfig:"USE_AUTH"`
	Password string `envconfig:"PASSWORD"`
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

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}

func GetBaseUrl() string {
	url := fmt.Sprintf("%s://%s", ConfigHandy.HttpProtocol, ConfigHandy.RefererHost)

	if ConfigHandy.HTTPPortUrl != 80 {
		url = fmt.Sprintf("%s:%d", url, ConfigHandy.HTTPPortUrl)
	}

	return url
}

func LoadConfigEnvApp() *Config {

	// Determine which .env file to load
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")

	// Load environment variables from .env file
	var err error
	switch env {
	case "production":
		err = godotenv.Load(".env.production")
	case "test":
		err = godotenv.Load("../../../../../../.env.test")
	case "local":
		fallthrough
	default:
		err = godotenv.Load(".env.local")
	}

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
	// Default Graylog file based on the environment variable
	defaultConfigFile := "config/config-local.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("config/config-%s.yaml", env)
	}

	// Load Master Graylog File
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
