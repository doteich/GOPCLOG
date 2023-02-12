package setup

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ClientConfig   ClientConfig `mapstructure:"opcConfig"`
	Nodes          []NodeObject `mapstructure:"selectedTags"`
	LoggerConfig   LoggerConfig `mapstructure:"methodConfig"`
	ExporterConfig Exporters    `mapstructure:"exporters"`
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
	Interval int    `mapstructure:"subInterval"`
	Name     string `mapstructure:"name"`
}
type Exporters struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Rest       RestConfig       `mapstructure:"rest"`
	Websockets WebsocketConfig  `mapstructure:"websockets"`
}

type PrometheusConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type RestConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"targetURL"`
}

type WebsocketConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

var PubConfig Config

func SetConfig() *Config {
	viper.AddConfigPath("/etc/config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	// Mount Secret to the application from env

	viper.Unmarshal(&PubConfig)

	if PubConfig.ClientConfig.AuthType != "Anonymous" {
		user := os.Getenv("OPCUA_USERNAME")
		pw := os.Getenv("OPCUA_PASSWORD")

		PubConfig.ClientConfig.Username = user
		PubConfig.ClientConfig.Password = pw
	}

	return &PubConfig
}
