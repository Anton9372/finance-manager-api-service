package config

import (
	"finance-manager-api-service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	JWT struct {
		Secret string `yaml:"secret" env-required:"true"`
	}
	HTTP struct {
		IP   string `yaml:"ip"`
		Port int    `yaml:"port"`
		CORS struct {
			AllowedMethods   []string `yaml:"allowed_methods"`
			AllowedOrigins   []string `yaml:"allowed_origins"`
			AllowCredentials bool     `yaml:"allow_credentials"`
			AllowedHeaders   []string `yaml:"allowed_headers"`
			ExposedHeaders   []string `yaml:"exposed_headers"`
		} `yaml:"cors"`
	} `yaml:"http"`
	UserService struct {
		HttpUrl         string `yaml:"http_url" env-required:"true"`
		GrpcUrl         string `yaml:"grpc_url" env-required:"true"`
		ConnectWithGRPC bool   `yaml:"connect_with_grpc"`
	} `yaml:"user_service" env-required:"true"`
	OperationService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"operation_service" env-required:"true"`
	StatsService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"stats_service" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/local.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
