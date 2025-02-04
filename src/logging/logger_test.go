package logging

import (
	"bytes"
	"strconv"
	//"fmt"
	//	"context"
	//	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"

	"sync"
	"testing"
	// "time"
)

// Custom syncWriter for controlled output capture (used for capturing log prints)
type syncWriter struct {
	Output *bytes.Buffer
}

// Replaceable stdout writer
var stdOut = &syncWriter{Output: new(bytes.Buffer)}

// Helper to capture output printed to stdout using os.Pipe()
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout // Backup the original stdout
	defer func() {
		os.Stdout = stdout // Restore original stdout after test
	}()

	// Redirect stdout to the pipe
	os.Stdout = w

	// Run the function that produces output
	f()
	w.Close()

	// Read the output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestLoggerLevels(t *testing.T) {
	handler := NewHandler(nil)

	output := captureOutput(func() {
		log := slog.New(handler)

		// Log messages at different levels
		log.Debug("Debugging")
		log.Info("Application started")
		log.Warn("Warning issued")
		log.Error("Error occurred")
	})

	// Log captured output for inspection
	t.Logf("Captured output:\n%s", output)

	// Verify that all log messages are present
	if !strings.Contains(output, "Debugging") {
		t.Errorf("Expected 'Debugging' message not found")
	}
	if !strings.Contains(output, "Application started") {
		t.Errorf("Expected 'Application started' message not found")
	}
	if !strings.Contains(output, "Warning issued") {
		t.Errorf("Expected 'Warning issued' message not found")
	}
	if !strings.Contains(output, "Error occurred") {
		t.Errorf("Expected 'Error occurred' message not found")
	}
}

// Test attributes passed to the logger
func TestLoggerWithAttrs(t *testing.T) {
	handler := NewHandler(nil)

	output := captureOutput(func() {
		log := slog.New(handler.WithAttrs(
			[]slog.Attr{
				slog.String("env", "production"),
				slog.Int("version", 2),
			},
		))
		log.Info("Log with attributes", slog.String("module", "auth"))
	})

	// Log captured output for inspection
	t.Logf("Captured output:\n%s", output)

	// Verify attribute content in the output
	if !strings.Contains(output, `"env": "production"`) {
		t.Errorf("Expected 'env' attribute not found")
	}
	if !strings.Contains(output, `"version": 2`) {
		t.Errorf("Expected 'version' attribute not found")
	}
	if !strings.Contains(output, `"module": "auth"`) {
		t.Errorf("Expected 'module' attribute not found")
	}
}

// Test log grouping
func TestLoggerWithGroup(t *testing.T) {
	handler := NewHandler(nil)
	output := captureOutput(func() {
		log := slog.New(handler.WithGroup("request"))
		log.Info("Grouped log", slog.String("path", "/api/v1/resource"), slog.Int("status", 200))
	})

	// Verify grouping
	if !strings.Contains(output, `"request":`) {
		t.Errorf("Expected 'request' group not found")
	}
	if !strings.Contains(output, `"path": "/api/v1/resource"`) {
		t.Errorf("Expected 'path' attribute not found in grouped log")
	}
	if !strings.Contains(output, `"status": 200`) {
		t.Errorf("Expected 'status' attribute not found in grouped log")
	}
	t.Log("Group logging test passed.")
}

// Test log concurrency safety (multiple goroutines logging simultaneously)
func TestLoggerConcurrency(t *testing.T) {
	handler := NewHandler(nil)
	var wg sync.WaitGroup

	output := captureOutput(func() {
		// Launch multiple goroutines to log concurrently
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				log := slog.New(handler)
				log.Info("Concurrent log", slog.Int("goroutine", i))
			}(i)
		}
		wg.Wait()
	})

	// Check that logs from multiple goroutines are present
	for i := 0; i < 10; i++ {
		if !strings.Contains(output, `"goroutine": `+strconv.Itoa(i)) {
			t.Errorf("Log entry from goroutine %d not found", i)
		}
	}
	t.Log("Concurrency test passed.")
}

// Test that custom log messages with errors are formatted correctly
func TestLoggerErrorHandling(t *testing.T) {
	handler := NewHandler(nil)
	output := captureOutput(func() {
		log := slog.New(handler)
		log.Error(
            "Operation failed",
            slog.String("reason", "network timeout"), 
            slog.String("error", errors.New("timeout").Error()),
        )
	})

	// Validate that error-related content appears correctly
	if !strings.Contains(output, "Operation failed") {
		t.Errorf("Expected 'Operation failed' message not found")
	}
	if !strings.Contains(output, `"reason": "network timeout"`) {
		t.Errorf("Expected 'reason' attribute not found")
	}
	if !strings.Contains(output, `"error": "timeout"`) {
		t.Errorf("Expected 'error' attribute not found")
	}
	t.Log("Error handling test passed.")
}
