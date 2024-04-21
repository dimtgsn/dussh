package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"dev"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
	Logger      `yaml:"logger" env-required:"true"`
	DB          `yaml:"database" env-required:"true"`
	Redis       `yaml:"redis" env-required:"true"`
	Auth        `yaml:"auth" env-required:"true"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
	User            string        `yaml:"user" env-required:"true"`
	Password        string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type DB struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
}

type Redis struct {
	Addr     string `yaml:"addr" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

type Auth struct {
	SecretKey       string        `yaml:"secret_key" env-required:"true" env:"AUTH_SECRET_KEY"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-required:"true"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-required:"true"`
}

type Logger struct {
	Level    string `yaml:"log_level" env-default:"debug"`
	Encoding string `yaml:"encoding" env-default:"json"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
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
