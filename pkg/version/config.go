package version

import (
	"log"

	"github.com/spf13/viper"
)

// Config type
type Config interface {
	GetString(string) string
}

// GetConfig gets new viper config instance populated with config data
func GetConfig(filename string, path string) (Config, error) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	viper.BindEnv("redis.address", "REDIS_ADDRESS")

	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file %s not found, falling back to default and ENV\n", filename)
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	return viper.GetViper(), nil
}
