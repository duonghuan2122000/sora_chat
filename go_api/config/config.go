// config.go

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HostName     string   `mapstructure:"HOST_NAME"`
	AppPort      string   `mapstructure:"APP_PORT"`
	TrustProxies []string `mapstructure:"TRUST_PROXIES"`
	MongoURI     string   `mapstructure:"MONGO_URI"`
	MongoDB      string   `mapstructure:"MONGO_DB"`
	RedisConn    string   `mapstructure:"REDIS_CONN"`
	RedisPass    string   `mapstructure:"REDIS_PASS"`
}

var AppConfig Config

func LoadConfig(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv() // ưu tiên ENV thực nếu có

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		fmt.Printf("Unable to decode config into struct: %v", err)
	}
}
