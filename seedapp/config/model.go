package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var osExit = os.Exit

type Config struct {
	App                   string         `mapstructure:"app"`
	AppVer                string         `mapstructure:"app_version"`
	Env                   string         `mapstructure:"env"`
	Http                  HttpConfig     `mapstructure:"http"`
	Log                   LogConfig      `mapstructure:"log"`
	Database              DatabaseConfig `mapstructure:"database"`
	Toggle                ToggleConfig   `mapstructure:"toggle"`
	BasicAuths            string         `mapstructure:"basic_auths"`
	FlagSeedDemoDate      time.Time      `mapstructure:"flag_seed_demo_date"`
	MigrateTenantStartup  bool           `mapstructure:"migrate_tenant_startup"`
	MigrateTenantPerBatch int            `mapstructure:"migrate_tenant_per_batch"`
	Kokai                 Kokai          `mapstructure:"kokai"`
	Redis                 RedisConfig    `mapstructure:"redis"`
	SchedulerTime         SchedulerTime  `mapstructure:"scheduler_time"`
}

func (c *Config) GetAppName() string {
	return c.App + "-" + c.Env
}

type HttpConfig struct {
	Port         int `mapstructure:"port"`
	WriteTimeout int `mapstructure:"write_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
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
	Host            string `mapstructure:"host"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	Port            string `mapstructure:"port"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
}

type ToggleConfig struct {
	AppName string `mapstructure:"app_name"`
	URL     string `mapstructure:"url"`
	Token   string `mapstructure:"token"`
}

type Kokai struct {
	URL      string        `mapstructure:"url"`
	User     string        `mapstructure:"user"`
	Password string        `mapstructure:"password"`
	Timeout  time.Duration `mapstructure:"timeout"`
}
type RedisConfig struct {
	Mode     string `mapstructure:"mode"`
	Address  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}
type SchedulerTime struct {
	PurchaseInvoiceRoutine string `mapstructure:"purchase_invoice_routine"`
	SalesInvoiceRoutine    string `mapstructure:"sales_invoice_routine"`
	JournalRoutine         string `mapstructure:"journal_routine"`
	VoucherRoutine         string `mapstructure:"voucher_routine"`
}

func (c *Config) LoadConfig(path string) {
	viper.AddConfigPath(".")
	viper.SetConfigName(path)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		osExit(1)
	}

	// Set up environment variable mappings using struct tags
	viper.AutomaticEnv()

	// Set up a replacer to match the struct tags
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		osExit(1)
	}
}
