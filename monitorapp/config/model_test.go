package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	testConfigDir string
	originalEnv   map[string]string
}

func (suite *ConfigTestSuite) SetupTest() {
	// Save original environment variables
	suite.originalEnv = make(map[string]string)

	// Reset viper for each test
	viper.Reset()

	// Create a temporary directory for test config files
	suite.testConfigDir = suite.T().TempDir()
}

func (suite *ConfigTestSuite) TearDownTest() {
	// Restore original environment variables
	for key, value := range suite.originalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}

	// Clean up environment variables we set
	envVars := []string{
		"APP", "APP_VERSION", "ENV",
		"HTTP_PORT", "HTTP_WRITE_TIMEOUT", "HTTP_READ_TIMEOUT",
		"LOG_FILE_LOCATION", "LOG_FILE_TDR_LOCATION", "LOG_FILE_MAX_SIZE",
		"LOG_FILE_MAX_BACKUP", "LOG_FILE_MAX_AGE", "LOG_STDOUT",
		"DATABASE_HOST", "DATABASE_USER", "DATABASE_PASSWORD", "DATABASE_NAME",
		"DATABASE_PORT", "DATABASE_SSL_MODE", "DATABASE_MAX_IDLE_CONN",
		"DATABASE_CONN_MAX_LIFETIME", "DATABASE_MAX_OPEN_CONN",
		"REDIS_MODE", "REDIS_ADDRESS", "REDIS_PORT", "REDIS_PASSWORD",
		"TOGGLE_APP_NAME", "TOGGLE_URL", "TOGGLE_TOKEN",
	}
	for _, key := range envVars {
		os.Unsetenv(key)
	}
}

func (suite *ConfigTestSuite) createTestConfigFile(filename string, content string) string {
	filePath := filepath.Join(suite.testConfigDir, filename+".yaml")
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

func (suite *ConfigTestSuite) TestLoadConfig_Success() {
	configContent := `
app: test-app
app_version: v1.0.0
env: test

http:
  port: 8080
  write_timeout: 60
  read_timeout: 60

log:
  file_location: "/var/log/test"
  file_tdr_location: "/var/log/test/tdr"
  file_max_size: 100
  file_max_backup: 5
  file_max_age: 7
  stdout: false

database:
  host: localhost
  user: testuser
  password: testpass
  name: testdb
  port: "5432"
  ssl_mode: require
  max_idle_conn: 5
  conn_max_lifetime: 2h
  max_open_conn: 50

redis:
  mode: single
  address: 127.0.0.1
  port: 6379
  password: redispass

toggle:
  app_name: test-app
  url: https://test.example.com/api/
  token: test-token
`

	suite.createTestConfigFile("test-config", configContent)

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("test-config")

	// Verify top-level fields
	suite.Equal("test-app", config.App)
	suite.Equal("v1.0.0", config.AppVer)
	suite.Equal("test", config.Env)

	// Verify HTTP config
	suite.Equal(8080, config.Http.Port)
	suite.Equal(60, config.Http.WriteTimeout)
	suite.Equal(60, config.Http.ReadTimeout)

	// Verify Log config
	suite.Equal("/var/log/test", config.Log.FileLocation)
	suite.Equal("/var/log/test/tdr", config.Log.FileTDRLocation)
	suite.Equal(100, config.Log.FileMaxSize)
	suite.Equal(5, config.Log.FileMaxBackup)
	suite.Equal(7, config.Log.FileMaxAge)
	suite.Equal(false, config.Log.Stdout)

	// Verify Database config
	suite.Equal("localhost", config.Database.Host)
	suite.Equal("testuser", config.Database.Username)
	suite.Equal("testpass", config.Database.Password)
	suite.Equal("testdb", config.Database.DBName)
	suite.Equal("5432", config.Database.Port)
	suite.Equal("require", config.Database.SSLMode)
	suite.Equal(5, config.Database.MaxIdleConn)
	suite.Equal(2*time.Hour, config.Database.ConnMaxLifetime)
	suite.Equal(50, config.Database.MaxOpenConn)

	// Verify Redis config
	suite.Equal("single", config.Redis.Mode)
	suite.Equal("127.0.0.1", config.Redis.Address)
	suite.Equal(6379, config.Redis.Port)
	suite.Equal("redispass", config.Redis.Password)

	// Verify Toggle config
	suite.Equal("test-app", config.Toggle.AppName)
	suite.Equal("https://test.example.com/api/", config.Toggle.URL)
	suite.Equal("test-token", config.Toggle.Token)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithEnvironmentVariables() {
	configContent := `
app: original-app
app_version: v0.0.1
env: development

http:
  port: 8000
  write_timeout: 30
  read_timeout: 30

log:
  file_location: "logs"
  file_tdr_location: "logs"
  file_max_size: 20
  file_max_backup: 10
  file_max_age: 30
  stdout: true

database:
  host: 127.0.0.1
  user: postgres
  password: postgres
  name: pakuningratan
  port: "5432"
  ssl_mode: disable
  max_idle_conn: 10
  conn_max_lifetime: 1h
  max_open_conn: 100

redis:
  mode: single
  address: 127.0.0.1
  port: 6379
  password: ""

toggle:
  app_name: original-toggle-app
  url: ""
  token: ""
`

	suite.createTestConfigFile("test-env-config", configContent)

	// Set environment variables to override config
	os.Setenv("APP", "env-app")
	os.Setenv("APP_VERSION", "v2.0.0")
	os.Setenv("HTTP_PORT", "9000")
	os.Setenv("LOG_STDOUT", "false")
	os.Setenv("DATABASE_HOST", "env-db-host")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("TOGGLE_APP_NAME", "env-toggle-app")

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("test-env-config")

	// Verify environment variable overrides
	suite.Equal("env-app", config.App)
	suite.Equal("v2.0.0", config.AppVer)
	suite.Equal(9000, config.Http.Port)
	suite.Equal(false, config.Log.Stdout)
	suite.Equal("env-db-host", config.Database.Host)
	suite.Equal(6380, config.Redis.Port)
	suite.Equal("env-toggle-app", config.Toggle.AppName)

	// Verify non-overridden values remain from config file
	suite.Equal(30, config.Http.WriteTimeout)
	suite.Equal(30, config.Http.ReadTimeout)
	suite.Equal("single", config.Redis.Mode)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithMinimalConfig() {
	configContent := `
app: minimal-app
app_version: v0.1.0
env: production
`

	suite.createTestConfigFile("minimal-config", configContent)

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("minimal-config")

	// Verify minimal fields are loaded
	suite.Equal("minimal-app", config.App)
	suite.Equal("v0.1.0", config.AppVer)
	suite.Equal("production", config.Env)

	// Verify nested structs have zero values
	suite.Equal(0, config.Http.Port)
	suite.Equal("", config.Database.Host)
	suite.Equal("", config.Redis.Mode)
	suite.Equal("", config.Toggle.AppName)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithDurationValues() {
	configContent := `
app: duration-test
app_version: v1.0.0
env: test

database:
  conn_max_lifetime: 30m
`

	suite.createTestConfigFile("duration-config", configContent)

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("duration-config")

	// Verify duration is parsed correctly
	suite.Equal(30*time.Minute, config.Database.ConnMaxLifetime)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithBooleanValues() {
	configContent := `
app: bool-test
app_version: v1.0.0
env: test

log:
  stdout: false
`

	suite.createTestConfigFile("bool-config", configContent)

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("bool-config")

	// Verify boolean is parsed correctly
	suite.Equal(false, config.Log.Stdout)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithNumericValues() {
	configContent := `
app: numeric-test
app_version: v1.0.0
env: test

http:
  port: 3000
  write_timeout: 45
  read_timeout: 45

log:
  file_max_size: 500
  file_max_backup: 20
  file_max_age: 60

database:
  max_idle_conn: 15
  max_open_conn: 200

redis:
  port: 7000
`

	suite.createTestConfigFile("numeric-config", configContent)

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("numeric-config")

	// Verify numeric values are parsed correctly
	suite.Equal(3000, config.Http.Port)
	suite.Equal(45, config.Http.WriteTimeout)
	suite.Equal(45, config.Http.ReadTimeout)
	suite.Equal(500, config.Log.FileMaxSize)
	suite.Equal(20, config.Log.FileMaxBackup)
	suite.Equal(60, config.Log.FileMaxAge)
	suite.Equal(15, config.Database.MaxIdleConn)
	suite.Equal(200, config.Database.MaxOpenConn)
	suite.Equal(7000, config.Redis.Port)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithStringValues() {
	configContent := `
app: string-test
app_version: v1.0.0
env: test

log:
  file_location: "/custom/log/path"
  file_tdr_location: "/custom/log/tdr"

database:
  host: "custom-host"
  user: "custom-user"
  password: "custom-password"
  name: "custom-db"
  port: "3306"
  ssl_mode: "verify-full"

redis:
  mode: "cluster"
  address: "redis.example.com"
  password: "secret-password"

toggle:
  app_name: "custom-app"
  url: "https://custom.example.com/api/"
  token: "custom-token-12345"
`

	suite.createTestConfigFile("string-config", configContent)

	// Change to test directory so viper can find the config
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.testConfigDir)

	config := &Config{}
	config.LoadConfig("string-config")

	// Verify string values are parsed correctly
	suite.Equal("/custom/log/path", config.Log.FileLocation)
	suite.Equal("/custom/log/tdr", config.Log.FileTDRLocation)
	suite.Equal("custom-host", config.Database.Host)
	suite.Equal("custom-user", config.Database.Username)
	suite.Equal("custom-password", config.Database.Password)
	suite.Equal("custom-db", config.Database.DBName)
	suite.Equal("3306", config.Database.Port)
	suite.Equal("verify-full", config.Database.SSLMode)
	suite.Equal("cluster", config.Redis.Mode)
	suite.Equal("redis.example.com", config.Redis.Address)
	suite.Equal("secret-password", config.Redis.Password)
	suite.Equal("custom-app", config.Toggle.AppName)
	suite.Equal("https://custom.example.com/api/", config.Toggle.URL)
	suite.Equal("custom-token-12345", config.Toggle.Token)
}

func (suite *ConfigTestSuite) TestConfigStruct_ZeroValues() {
	// Test that structs have proper zero values
	config := &Config{}

	suite.Equal("", config.App)
	suite.Equal("", config.AppVer)
	suite.Equal("", config.Env)
	suite.Equal(0, config.Http.Port)
	suite.Equal("", config.Log.FileLocation)
	suite.Equal("", config.Database.Host)
	suite.Equal("", config.Redis.Mode)
	suite.Equal("", config.Toggle.AppName)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
