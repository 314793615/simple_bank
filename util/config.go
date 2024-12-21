package util

import (
	"fmt"

	"github.com/spf13/viper"
)


type Config struct {
	DBSource string
	DBDriver string
	Address string
}

func NewConfig(path string) (*Config, error ){
	path = ToAbsPath(path)
	config := &Config{}
	vip := viper.New()
	vip.SetConfigName("app")
	vip.AddConfigPath(path)

	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil{
		return nil, fmt.Errorf("failed to read in config of config file: %v", err)
	}
	err = vip.Unmarshal(&config)
	if err != nil{
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}
	return config, nil

}