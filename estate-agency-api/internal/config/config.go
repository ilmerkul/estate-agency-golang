package config

import (
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"gilab.com/estate-agency-api/internal/storage/database/mysql"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                 string `yaml:"env" env:"ENV" env-default:"debug"`
	mysql.StorageConfig `yaml:"storage"`
	HTTPServerConfig    `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
	TimeClose   time.Duration `yaml:"time_close" env:"HTTP_SERVER_TIME_CLOSE" env-default:"10s"`
	User        string        `yaml:"user" env:"HTTP_SERVER_USER" env-required:"true"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true" `
}

var once sync.Once

func MustLoad() *Config {
	var cfg *Config

	once.Do(
		func() {
			configPath := fetchConfigPath()
			if configPath == "" {
				panic("config path is empty")
			}

			cfg = MustLoadPath(configPath)
		},
	)

	return cfg
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		help, _ := cleanenv.GetDescription(cfg, nil)
		slog.Info(help)
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "../../config/config.yml", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
