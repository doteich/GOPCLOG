package setup

import "github.com/spf13/viper"

type Config struct {
	ClientConfig ClientConfig `mapstructure:"opcConfig"`
	Nodes        []NodeObject `mapstructure:"selectedTags"`
	LoggerConfig LoggerConfig `mapstructure:"methodConfig"`
}

type ClientConfig struct {
	Url            string `mapstructure:"url"`
	SecurityPolicy string `mapstructure:"securityPolicy"`
	SecurityMode   string `mapstructure:"securityMode"`
	AuthType       string `mapstructure:"authType"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	Node           string `mapstructure:"node"`
}
type NodeObject struct {
	NodeId   string `mapstructure:"nodeId"`
	NodeName string `mapstructure:"name"`
}
type LoggerConfig struct {
	Interval int    `mapstructure:"subInterval"`
	Name     string `mapstructure:"name"`
}

func SetConfig() Config {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	var c Config

	viper.Unmarshal(&c)
	return c
}
