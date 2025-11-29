package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// Init initializes a global structured logger which writes JSON logs to both
// stdout (so docker logs show them) and to a rotating file for persistence.
// logFile is the path to the log file (e.g. "logs/app.log").
func Init(logFile string, level logrus.Level) {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0o755); err != nil {
		// If we cannot create logs dir, still proceed and write to stdout
		// but we can't return an error here because logger init should be best-effort.
		// fallthrough
	}

	rotate := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Write to both stdout (docker logs) and the rotating file
	mw := io.MultiWriter(os.Stdout, rotate)
	l.SetOutput(mw)
	l.SetLevel(level)

	Log = l
}

// Helper wrappers
func Info(args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Infof(format, args...)
}

func Error(args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	if Log == nil {
		return
	}
	Log.Errorf(format, args...)
}

// LogrusLevel returns the default logging level used by the project.
// Exported so callers can pass a level without depending on logrus directly.
func LogrusLevel() logrus.Level {
	return logrus.InfoLevel
}
