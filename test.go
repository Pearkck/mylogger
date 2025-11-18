package loggerLib

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLogger_InfoAndError(t *testing.T) {
	ctx := context.Background()

	// Create temp dir for logs
	tempDir, err := ioutil.TempDir("", "logtest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up after test

	// Initialize logger
	logger, err := NewLogger("testservice")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test info log
	logger.Info(ctx, "Test info message")
	logger.Infof(ctx, "Formatted info %d", 123)

	// Test error log
	logger.Error(ctx, "Test error message")
	logger.Errorf(ctx, "Formatted error %s", "oops")

	// Check log files exist
	infoPath := filepath.Join(tempDir, "info.log")
	errorPath := filepath.Join(tempDir, "error.log")

	if _, err := os.Stat(infoPath); os.IsNotExist(err) {
		t.Errorf("info.log not created")
	}
	if _, err := os.Stat(errorPath); os.IsNotExist(err) {
		t.Errorf("error.log not created")
	}

	// Optionally, check log content
	infoContent, _ := os.ReadFile(infoPath)
	errorContent, _ := os.ReadFile(errorPath)

	if len(infoContent) == 0 {
		t.Errorf("info.log is empty")
	}
	if len(errorContent) == 0 {
		t.Errorf("error.log is empty")
	}
}
