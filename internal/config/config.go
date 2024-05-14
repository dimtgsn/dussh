package config

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Auth       `yaml:"auth" env-required:"true"`
	Env        string `yaml:"env" env-default:"dev"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	Logger     `yaml:"logger" env-required:"true"`
	DB         `yaml:"database" env-required:"true"`
	Redis      `yaml:"redis" env-required:"true"`
	RabbitMQ   `yaml:"rabbit_mq" env-required:"true"`
	Notify     `yaml:"notify" env-required:"true"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

type DB struct {
	Host     string `yaml:"host" env:"DATABASE_HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"DATABASE_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DATABASE_USER" env-required:"true"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD" env-required:"true"`
	DBName   string `yaml:"db_name" env:"DATABASE_NAME" env-required:"true"`
}

type Redis struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-required:"true"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-required:"true"`
}

type RabbitMQ struct {
	Port                  int    `yaml:"port" env:"RABBITMQ_PORT" env-default:"5672"`
	User                  string `yaml:"user" env:"RABBITMQ_USER" env-required:"true"`
	Host                  string `yaml:"host" env:"RABBITMQ_HOST" env-required:"true"`
	Password              string `yaml:"password" env:"RABBITMQ_PASSWORD" env-required:"true"`
	NotificationPublisher `yaml:"notification_publisher"`
	NotificationConsumer  `yaml:"notification_consumer"`
}

type NotificationPublisher struct {
	Exchange string `yaml:"exchange"`
}

type NotificationConsumer struct {
	Name  string `yaml:"name"`
	Queue string `yaml:"queue"`
}

type Notify struct {
	EmailProvider `yaml:"email_provider"`
}

type EmailProvider struct {
	Port      int    `yaml:"port" env:"NOTIFY_EMAIL_PROVIDER_PORT" env-default:"2525"`
	Host      string `yaml:"host" env:"NOTIFY_EMAIL_PROVIDER_HOST" env-required:"true"`
	From      string `yaml:"from" env-default:"dussh@school.com"`
	Username  string `yaml:"username" env:"NOTIFY_EMAIL_PROVIDER_USERNAME" env-required:"true"`
	Password  string `yaml:"password" env:"NOTIFY_EMAIL_PROVIDER_PASSWORD" env-required:"true"`
	TLSEnable bool   `yaml:"tls_enable" env:"NOTIFY_EMAIL_PROVIDER_TLS" env-default:"false"`
}

type Auth struct {
	SecretKey       string        `yaml:"secret_key" env:"AUTH_SECRET_KEY" env-required:"true"`
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
	// load .env file
	if err := godotenv.Load(); err != nil {
		panic("error loading .env file")
	}

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
