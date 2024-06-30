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
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	UserService struct {
		URL string `yaml:"url" env-required:"true"`
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
