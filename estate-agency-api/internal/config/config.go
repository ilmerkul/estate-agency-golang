package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string `yaml:"env" env:"ENV" env-default:"local"`
	StorageConfig    `yaml:"storage"`
	HTTPServerConfig `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
	User        string        `yaml:"user" env:"HTTP_SERVER_USER" env-required:"true"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true" `
}

type StorageConfig struct {
	StoragePath     string `yaml:"storage_path" env:"STORAGE_PATH"`
	ConfigDSN       `yaml:"dsn"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"0"`
	MaxIdleConns    int `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"50"`
	MaxOpenConns    int `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"50"`
}

type ConfigDSN struct {
	User     string `yaml:"user" env:"STORAGE_DSN_USER"`
	Password string `yaml:"password" env:"STORAGE_DSN_PASSWORD"`
	Protocol string `yaml:"protocol" env:"STORAGE_DSN_PROTOCOL"`
	Host     string `yaml:"host" env:"STORAGE_DSN_HOST"`
	Port     string `yaml:"port" env:"STORAGE_DSN_PORT"`
	NameDB   string `yaml:"name_db" env:"STORAGE_DSN_NAME_DB"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Println("CONFIG_PATH is not set")
		configPath = "C:\\Users\\pynex\\Projects\\estate-agency-api\\config\\config.yml"
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("Read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("C:\\Users\\pynex\\Projects\\estate-agency-api\\internal\\config\\config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
