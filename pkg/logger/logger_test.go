package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Critical test: verification of SetOutput
	// If this fails, other tests that rely on capturing output will likely fail or give false negatives.
	success := t.Run("TestSetOutput", testSetOutput)
	if !success {
		t.Fatal("Critical test TestSetOutput failed. Aborting remaining tests.")
	}

	// Run other tests
	t.Run("TestSetFormat", testSetFormat)
	t.Run("TestPartialUpdate", testPartialUpdate)
	t.Run("TestCleanSourcePath", testCleanSourcePath)
	t.Run("TestInvalidLevel", testInvalidLevel)
}

// --------------------------------
// Critical Tests
// --------------------------------

// If this fails, other tests that rely on capturing output will likely fail or give false negatives.
func testSetOutput(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	slog.Info("output test")

	if !strings.Contains(buf.String(), "output test") {
		t.Errorf("Expected output to contain 'output test', got '%s'", buf.String())
	}
}

// --------------------------------
// Tests
// --------------------------------
func testSetFormat(t *testing.T) {
	// Redirect output to a buffer
	var buf bytes.Buffer
	SetOutput(&buf)

	// Reset configuration for test
	color := false
	level := "DEBUG"
	err := SetConfig(Config{
		Colors: &color, // Ensure JSON output for parsing
		Level:  &level,
	})
	if err != nil {
		t.Fatalf("SetFormat failed: %v", err)
	}

	// Log something
	buf.Reset()
	checkCorrectLogContent(t, &buf, slog.Info, "test message", "key", "value", "INFO")
}

func testPartialUpdate(t *testing.T) {
	// Setup initial state
	var buf bytes.Buffer
	SetOutput(&buf)

	initialColor := false
	initialLevel := "INFO"
	SetConfig(Config{Colors: &initialColor, Level: &initialLevel})

	// Partial update: Change level only
	newLevel := "DEBUG"
	err := SetConfig(Config{Level: &newLevel})
	if err != nil {
		t.Fatalf("SetFormat failed: %v", err)
	}

	// Verify level changed
	checkCorrectLogContent(t, &buf, slog.Debug, "test message debug", "key", "value", newLevel)

	// Reset buffer
	buf.Reset()

	// Verify colors (format) didn't revert to default (which might be different)
	// Since we are checking JSON output (colors=false), we expect JSON.
	checkCorrectLogContent(t, &buf, slog.Info, "json check", "key", "value", "INFO")
	if !strings.HasPrefix(strings.TrimSpace(buf.String()), "{") {
		t.Errorf("Expected JSON output, got: %s", buf.String())
	}
}

func testCleanSourcePath(t *testing.T) {
	// This function is internal and used by the handler.
	// We can test it by checking if the source attribute in logs is just the filename.

	var buf bytes.Buffer
	SetOutput(&buf)

	color := false
	level := "INFO"

	SetConfig(Config{Colors: &color, Level: &level})
	buf.Reset()
	slog.Info("test source")

	var logEntry map[string]any
	decoder := json.NewDecoder(&buf)
	if err := decoder.Decode(&logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// If it captured the configuration log instead of the "test source" log
	if logEntry["msg"] == "Logger configured" {
		if err := decoder.Decode(&logEntry); err != nil {
			t.Fatalf("Failed to parse second log output: %v", err)
		}
	}

	source, ok := logEntry["source"].(map[string]any)
	if !ok {
		// If AddSource is failing or not present, we might need to check how it's configured.
		// The logger config sets AddSource: true.
		t.Fatalf("Expected source field in log entry")
	}

	file, ok := source["file"].(string)
	if !ok {
		t.Fatalf("Expected file field in source")
	}

	// It should be a base filename, not a full path.
	if strings.Contains(file, "/") || strings.Contains(file, "\\") {
		t.Errorf("Expected base filename, got full path: %s", file)
	}
}

func testInvalidLevel(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	// Attempt to set an invalid level
	badLevel := "INVALID_LEVEL"
	err := SetConfig(Config{Level: &badLevel})

	// The current implementation logs an error but does not return an error for the caller
	// (it returns nil in SetFormat). It falls back to INFO.
	if err != nil {
		t.Logf("SetFormat returned error as expected (optional behavior): %v", err)
	}

	//Check error message
	var logEntry map[string]any
	decoder := json.NewDecoder(&buf)
	for {
		if err := decoder.Decode(&logEntry); err != nil {
			t.Fatalf("Failed to parse log output: %v", err)
		}
		if logEntry["msg"] == ErrInvalidLevel {
			break
		}
		if decoder.Buffered() == nil {
			t.Fatalf("Expected message '%s' not found in logs", ErrInvalidLevel)
		}
	}

	// Verify fallback to INFO (default fallback in code)
	slog.Info("should be visible")
	if buf.Len() == 0 {
		t.Errorf("Expected generic info log to work after invalid level")
	}

	buf.Reset()
	slog.Debug("should not be visible if fallback is INFO")
	// Note: If default logic changes, this might fail, but currently fallback is INFO.
	if buf.Len() > 0 {
		t.Errorf("Expected DEBUG log to be hidden after invalid level fallback to INFO")
	}
}

// --------------------------------
// Helper Functions
// --------------------------------
func checkCorrectLogContent(t *testing.T,
	buf *bytes.Buffer, logFun func(msg string, args ...any),
	msg string, key string, keyValue string, level string,
) {
	buf.Reset()
	logFun(msg, key, keyValue)

	// Parse JSON output
	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v\nOutput: %s", err, buf.String())
	}

	if logEntry["msg"] != msg {
		t.Errorf("Expected message '%s', got '%s'", msg, logEntry["msg"])
	}
	if logEntry["key"] != keyValue {
		t.Errorf("Expected key '%s', got '%s'", keyValue, logEntry["key"])
	}
	if logEntry["level"] != level {
		t.Errorf("Expected level '%s', got '%v'", level, logEntry["level"])
	}

}
