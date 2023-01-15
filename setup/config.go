package setup

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

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
	GenerateCert   bool   `mapstructure:"autoGenCert"`
}
type NodeObject struct {
	NodeId          string `mapstructure:"nodeId"`
	NodeName        string `mapstructure:"name"`
	DataTypeId      int    `mapstructure:"dataTypeId"`
	DataType        string `mapstructure:"dataType"`
	ExposeAsMetrics bool   `mapstructure:"exposeAsMetric"`
	MetricsType     string `mapstructure:"metricsType"`
}
type LoggerConfig struct {
	Interval       int    `mapstructure:"subInterval"`
	Name           string `mapstructure:"name"`
	TargetURL      string `mapstructure:"targetURL"`
	MetricsEnabled bool   `mapstructure:"metricsEnabled"`
	BackupEnabled  bool   `mapstructure:"backup"`
}

func SetConfig() Config {
	viper.AddConfigPath("/etc/config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	// Mount Secret to the application from env

	var c Config

	viper.Unmarshal(&c)

	if c.ClientConfig.AuthType != "Anonymous" {
		user := os.Getenv("OPCUA_USERNAME")
		pw := os.Getenv("OPCUA_PASSWORD")

		c.ClientConfig.Username = user
		c.ClientConfig.Password = pw
	}

	fmt.Printf("User:%s; PW: %s", c.ClientConfig.Username, c.ClientConfig.Password)

	return c
}
