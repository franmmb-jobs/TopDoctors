package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lmittmann/tint"
)

type Config struct {
	Colors *bool
	Level  *string
	Output io.Writer
}
type internalConfig struct {
	colors bool
	level  slog.Level
	output io.Writer
}

// Current logger config values
var loggerCfg internalConfig

var mu sync.Mutex

// Errors
const (
	ErrInvalidLevel = "invalid log level"
)

// ----------------------------------------------------------------
// Initialize the logger
// ----------------------------------------------------------------
func init() {
	// Set our custom logger as the global default
	// We don't need to worry about races in init() because
	// Go guarantees that it runs in a single thread at startup.
	loggerCfg = internalConfig{
		colors: true,
		level:  slog.LevelDebug,
		output: os.Stdout,
	}
	// Apply initial configuration
	applyLogger()
}

// ----------------------------------------------------------------
// External functions
// ----------------------------------------------------------------

// SetConfig changes the global logger at runtime
// An error is returned even if not used because future logger changes might require it
func SetConfig(cfg Config) error {
	// Lock to prevent race conditions
	mu.Lock()
	defer mu.Unlock()

	/////////////////////////////
	// Check the configuration
	if cfg.Level != nil {
		errUMLevel := loggerCfg.level.UnmarshalText([]byte(strings.ToUpper(*cfg.Level)))
		if errUMLevel != nil {
			// If the level is misspelled, we warn about the failure but continue to avoid blocking the program
			slog.Error(ErrInvalidLevel, "error_detail", errUMLevel)
		}
	}
	if cfg.Colors != nil {
		loggerCfg.colors = *cfg.Colors
	}

	// Manage the output
	if cfg.Output != nil {
		loggerCfg.output = cfg.Output
	}
	// If no output is defined (it's the first time or hasn't been set), use os.Stdout by default
	if loggerCfg.output == nil {
		loggerCfg.output = os.Stdout
	}

	applyLogger()

	return nil
}

// SetOutput changes the log destination at runtime
// This function is kept external to the configuration for several reasons:
// 1. It is a functionality used in tests
// 2. It is a functionality often used at runtime, such as changing the log storage file
// However, file saving is not implemented in this library as its goal is to be as simple and lightweight as possible
// and let the user decide where to store the logs
func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	loggerCfg.output = w
	applyLogger()
}

// ----------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------

// applyLogger applies the current configuration to the global logger
func applyLogger() {
	// Define the handler
	var handler slog.Handler

	// Create based on whether colors are supported
	if loggerCfg.colors {
		handler = tint.NewHandler(loggerCfg.output, &tint.Options{
			Level:       loggerCfg.level,
			AddSource:   true,
			TimeFormat:  time.Kitchen,
			ReplaceAttr: cleanSourcePath,
		})
	} else {
		// JSON format because it's more useful for production, where colors are generally not used
		handler = slog.NewJSONHandler(loggerCfg.output, &slog.HandlerOptions{
			Level:       loggerCfg.level,
			AddSource:   true,
			ReplaceAttr: cleanSourcePath,
		})
	}

	// Set the new logger as the global default
	slog.SetDefault(slog.New(handler))

	// Log the configuration
	slog.Info("Logger configured", "colors", loggerCfg.colors, "level", loggerCfg.level)
}

// ----------------------------------------------------------------
// Auxiliary functions
// ----------------------------------------------------------------

// cleanSourcePath is an auxiliary function to clean the file path in the log
func cleanSourcePath(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey && a.Value.Kind() == slog.KindAny {
		if source, ok := a.Value.Any().(*slog.Source); ok {
			source.File = filepath.Base(source.File)
		}
	}
	return a
}
