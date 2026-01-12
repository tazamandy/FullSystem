package logging

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the global logger instance
var Logger *zap.Logger

// currentLogFile tracks the current log file for daily rotation
var currentLogFile string

// InitLogger initializes the modern logging system with daily log rotation
func InitLogger() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Generate initial filename with current date
	currentDate := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("logs/ATENDIFY-%s.log", currentDate)
	currentLogFile = filename

	// Create the log file immediately
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	file.Close()

	// Configure lumberjack for file rotation
	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // MB
		MaxBackups: 30,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Text encoder for file (standard logs)
	fileEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// Console encoder (human-readable with colors)
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// Create file syncer with daily rotation wrapper
	fileSyncer := NewDailyRotatingWriter(fileWriter)

	// Create cores
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// Combine cores (logs go to both file and console)
	core := zapcore.NewTee(fileCore, consoleCore)

	// Create logger with caller information
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// DailyRotatingWriter wraps lumberjack to handle daily log rotation
type DailyRotatingWriter struct {
	writer   *lumberjack.Logger
	lastDate string
	mu       sync.Mutex
}

// NewDailyRotatingWriter creates a new daily rotating writer
func NewDailyRotatingWriter(writer *lumberjack.Logger) zapcore.WriteSyncer {
	currentDate := time.Now().Format("2006-01-02")
	return &DailyRotatingWriter{
		writer:   writer,
		lastDate: currentDate,
	}
}

// Write implements io.Writer and handles daily rotation
func (d *DailyRotatingWriter) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")

	// Check if date has changed, update filename if needed
	if currentDate != d.lastDate {
		d.lastDate = currentDate
		newFilename := fmt.Sprintf("logs/ATENDIFY-%s.log", currentDate)
		d.writer.Filename = newFilename
		currentLogFile = newFilename
	}

	return d.writer.Write(p)
}

// Sync flushes any buffered log entries to disk
func (d *DailyRotatingWriter) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// lumberjack doesn't have Sync, but writes are flushed automatically
	return nil
}

// customTimeEncoder formats time in a readable way for console
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// LogWithIP creates a logger with IP address field
func LogWithIP(ip string) *zap.Logger {
	return Logger.With(zap.String("client_ip", ip))
}

// LogRequest logs HTTP requests with IP
func LogRequest(ip, method, path string, status int, duration time.Duration) {
	Logger.Info("HTTP Request",
		zap.String("client_ip", ip),
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.Duration("duration", duration),
	)
}

// LogAuth logs authentication events
func LogAuth(ip, action, userID string, success bool) {
	level := zapcore.InfoLevel
	if !success {
		level = zapcore.WarnLevel
	}
	Logger.Log(level, "Authentication",
		zap.String("client_ip", ip),
		zap.String("action", action),
		zap.String("user_id", userID),
		zap.Bool("success", success),
	)
}

// LogSecurity logs security-related events
func LogSecurity(ip, event, details string) {
	Logger.Warn("Security Event",
		zap.String("client_ip", ip),
		zap.String("event", event),
		zap.String("details", details),
	)
}

// LogError logs errors with context
func LogError(ip string, err error, context string) {
	Logger.Error("Error",
		zap.String("client_ip", ip),
		zap.Error(err),
		zap.String("context", context),
	)
}
