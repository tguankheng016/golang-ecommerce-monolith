package logging

import (
	"context"
	"os"

	"github.com/tguankheng016/commerce-mono/pkg/environment"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.Logger
)

// InitLogger sets up a logger according to the given environment.
// It logs both to stdout and to a file named "go.log" in the "log" directory.
// The file is rotated daily and the maximum size is 100MB, with up to 3 backups.
// The logger's level is set to Info in production and Debug in development.
func InitLogger(env environment.Environment) *zap.Logger {
	if env == "" {
		env = environment.Development
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "../../log/go.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	ws := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		w,
	)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var core zapcore.Core

	if env.IsDevelopment() {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			ws,
			zap.DebugLevel,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			ws,
			zap.InfoLevel,
		)
	}

	Logger = zap.New(core)

	return Logger
}

// RunLogger sets up the fx lifecycle hooks for the zap logger.
// It writes an info message when the application starts and another when the
// application is stopping, and syncs the logger on stop.
func RunLogger(lc fx.Lifecycle, logger *zap.Logger, ctx context.Context) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("starting logger...")

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("close and syncing logger...")
			logger.Sync()
			return nil
		},
	})

	return nil
}
