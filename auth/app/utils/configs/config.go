package configs

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Environment string

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	Env      Environment   `yaml:"env"`
	TokenTTL time.Duration `yaml:"tokenTTL"`
	HTTP     HTTP          `yaml:"http"`
	GRPC     GRPC          `yaml:"grpc"`
	Cache    Cache         `yaml:"cache"`
	DB       Database      `yaml:"db"`
}

type HTTP struct {
	Port int `yaml:"port"`
}

type GRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Cache struct {
	TTL     time.Duration `yaml:"ttl"`
	MaxSize int           `yaml:"max-size"`
}

type Database struct {
	Port               int    `yaml:"port"`
	Host               string `yaml:"host"`
	DBName             string `yaml:"db-name"`
	Username           string `yaml:"username"`
	SSLMode            string `yaml:"ssl-mode"`
	MaxConnections     int    `yaml:"max-connections"`
	MaxIdleConnections int    `yaml:"max-idle-connections"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
