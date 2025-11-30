package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{
			name:     "FATAL level",
			level:    LevelFatal,
			expected: "FATAL",
		},
		{
			name:     "ERROR level",
			level:    LevelError,
			expected: "ERROR",
		},
		{
			name:     "WARNING level",
			level:    LevelWarning,
			expected: "WARNING",
		},
		{
			name:     "INFO level",
			level:    LevelInfo,
			expected: "INFO",
		},
		{
			name:     "DEBUG level",
			level:    LevelDebug,
			expected: "DEBUG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.level) != tt.expected {
				t.Errorf("expected level '%s', got '%s'", tt.expected, tt.level)
			}
		})
	}
}

func TestLog(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Log(LevelInfo, "TestFeature", "Test message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Error("expected output to contain [INFO]")
	}
	if !strings.Contains(output, "[TestFeature]") {
		t.Error("expected output to contain [TestFeature]")
	}
	if !strings.Contains(output, "Test message") {
		t.Error("expected output to contain 'Test message'")
	}
}

func TestFatal(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Fatal("FatalFeature", "Fatal error occurred")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[FATAL]") {
		t.Error("expected output to contain [FATAL]")
	}
	if !strings.Contains(output, "[FatalFeature]") {
		t.Error("expected output to contain [FatalFeature]")
	}
}

func TestError(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Error("ErrorFeature", "Error occurred")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[ERROR]") {
		t.Error("expected output to contain [ERROR]")
	}
	if !strings.Contains(output, "[ErrorFeature]") {
		t.Error("expected output to contain [ErrorFeature]")
	}
}

func TestWarning(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Warning("WarningFeature", "Warning message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[WARNING]") {
		t.Error("expected output to contain [WARNING]")
	}
	if !strings.Contains(output, "[WarningFeature]") {
		t.Error("expected output to contain [WarningFeature]")
	}
}

func TestInfo(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Info("InfoFeature", "Info message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Error("expected output to contain [INFO]")
	}
	if !strings.Contains(output, "[InfoFeature]") {
		t.Error("expected output to contain [InfoFeature]")
	}
}

func TestDebug(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Debug("DebugFeature", "Debug message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[DEBUG]") {
		t.Error("expected output to contain [DEBUG]")
	}
	if !strings.Contains(output, "[DebugFeature]") {
		t.Error("expected output to contain [DebugFeature]")
	}
}

func TestLogFormat(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Log(LevelInfo, "Feature", "Message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "-") {
		t.Error("expected output to contain date separator")
	}
	if !strings.Contains(output, ":") {
		t.Error("expected output to contain time separator")
	}
}
