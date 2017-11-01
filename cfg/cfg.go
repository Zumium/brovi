package cfg

import (
	"github.com/spf13/viper"
)

//Init initializes the config system
func Init() error {
	viper.SetEnvPrefix("brovi")
	viper.AutomaticEnv()
	return nil
}

//Port returns the port setting
func Port() int {
	return viper.GetInt("server_port")
}
