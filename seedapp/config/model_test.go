package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Helper function to clean up viper after each test
	cleanup := func() {
		viper.Reset()
		os.Clearenv() // Make sure we clear all env vars
	}

	t.Run("successful config load", func(t *testing.T) {
		defer cleanup()

		// Create a temporary config file
		configContent := []byte(`
app: "testapp"
app_version: "1.0.0"
env: "test"
http:
  port: 8080
  write_timeout: 15
  read_timeout: 15
log:
  file_location: "test.log"
  file_max_size: 100
  file_max_backup: 10
  file_max_age: 30
  stdout: true
database:
  host: "localhost"
  port: "5432"
`)
		tmpfile, err := os.CreateTemp(".", "config-*.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(configContent); err != nil {
			t.Fatal(err)
		}
		tmpfile.Close()

		// Set test environment variable
		os.Setenv("DATABASE_HOST", "testhost")

		// Test config loading
		cfg := &Config{}
		cfg.LoadConfig(tmpfile.Name()[:len(tmpfile.Name())-5]) // Remove .yaml extension

		// Verify loaded configuration
		assert.Equal(t, "testapp", cfg.App)
		assert.Equal(t, "1.0.0", cfg.AppVer)
		assert.Equal(t, 8080, cfg.Http.Port)
		assert.Equal(t, "testhost", cfg.Database.Host)
	})

	t.Run("missing config file", func(t *testing.T) {
		defer cleanup()

		// Capture stdout
		oldStdout := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		// Restore stdout when done
		defer func() {
			os.Stdout = oldStdout
		}()

		// Capture os.Exit call
		originalOsExit := osExit
		defer func() {
			osExit = originalOsExit
			if r := recover(); r == nil {
				t.Error("Expected panic but got none")
			}
		}()

		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
			panic("os.Exit called")
		}

		cfg := &Config{}
		cfg.LoadConfig("nonexistent")

		if !exitCalled {
			t.Error("os.Exit was not called")
		}
	})

	t.Run("environment variable fallback", func(t *testing.T) {
		defer cleanup()

		// Create a temporary config file with env variable reference
		configContent := []byte(`
database:
  host: "fallback-host"
`)
		tmpfile, err := os.CreateTemp(".", "config-*.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(configContent); err != nil {
			t.Fatal(err)
		}
		tmpfile.Close()

		// Test without environment variable (should use fallback)
		cfg := &Config{}
		cfg.LoadConfig(tmpfile.Name()[:len(tmpfile.Name())-5])
		assert.Equal(t, "fallback-host", cfg.Database.Host)

		// Clear Viper's cache and set environment variable
		viper.Reset()
		os.Setenv("DATABASE_HOST", "env-host")

		// Test with environment variable
		cfg = &Config{}
		cfg.LoadConfig(tmpfile.Name()[:len(tmpfile.Name())-5])
		assert.Equal(t, "env-host", cfg.Database.Host)
	})
}
