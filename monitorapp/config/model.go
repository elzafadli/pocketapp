package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App            string         `mapstructure:"app"`
	AppVer         string         `mapstructure:"app_version"`
	Env            string         `mapstructure:"env"`
	Http           HttpConfig     `mapstructure:"http"`
	Grpc           GrpcConfig     `mapstructure:"grpc"`
	Log            LogConfig      `mapstructure:"log"`
	Database       DatabaseConfig `mapstructure:"database"`
	Redis          RedisConfig    `mapstructure:"redis"`
	Toggle         ToggleConfig   `mapstructure:"toggle"`
	BasicAuths     string         `mapstructure:"basic_auths"`
	InternalApiKey string         `mapstructure:"internal_api_key"`
}

func (c *Config) GetAppName() string {
	return c.App + "-" + c.Env
}

type HttpConfig struct {
	Port         int `mapstructure:"port"`
	WriteTimeout int `mapstructure:"write_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
}

type GrpcConfig struct {
	Port int `mapstructure:"port"`
}

type LogConfig struct {
	FileLocation    string `mapstructure:"file_location"`
	FileTDRLocation string `mapstructure:"file_tdr_location"`
	FileMaxSize     int    `mapstructure:"file_max_size"`
	FileMaxBackup   int    `mapstructure:"file_max_backup"`
	FileMaxAge      int    `mapstructure:"file_max_age"`
	Stdout          bool   `mapstructure:"stdout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Username        string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"name"`
	Port            string        `mapstructure:"port"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxIdleConn     int           `mapstructure:"max_idle_conn"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
}

type RedisConfig struct {
	Mode     string `mapstructure:"mode"`
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type ToggleConfig struct {
	AppName string `mapstructure:"app_name"`
	URL     string `mapstructure:"url"`
	Token   string `mapstructure:"token"`
}

func (c *Config) LoadConfig(path string) {
	viper.AddConfigPath(".")
	viper.SetConfigName(path)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
}
