package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	// Test successful execution
	t.Run("successful execution", func(t *testing.T) {
		// This should not panic
		Execute()
	})

	t.Run("command error", func(t *testing.T) {
		// Temporarily modify rootCmd to simulate an error
		originalCmd := rootCmd
		defer func() {
			rootCmd = originalCmd // Restore original command
			// Recover from panic
			if r := recover(); r == nil {
				t.Error("Expected panic but got none")
			}
		}()

		// Create a command that will generate an error
		rootCmd = &cobra.Command{
			Use: "test",
			RunE: func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("test error")
			},
		}

		Execute() // This should panic
	})
}
