package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

func LevelFromString(value string) Level {
	switch value {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel
	case "dpanic":
		return DPanicLevel
	case "fatal":
		return FatalLevel
	default:
		return ErrorLevel
	}
}

const (
	DebugLevel  = Level(zapcore.DebugLevel)
	InfoLevel   = Level(zapcore.InfoLevel)
	WarnLevel   = Level(zapcore.WarnLevel)
	ErrorLevel  = Level(zapcore.ErrorLevel)
	PanicLevel  = Level(zapcore.PanicLevel)
	DPanicLevel = Level(zapcore.DPanicLevel)
	FatalLevel  = Level(zapcore.FatalLevel)
)

func newZapFileLogger(conf Config) (_ *zap.Logger, fileCloser func(), err error) {
	currentLevel := zapcore.Level(conf.Level)
	coreLoggers := make([]zapcore.Core, 0)
	encoderConfig := newZapEncodeConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	currentLevelPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= currentLevel
	})
	if len(conf.Path) > 0 {
		var fileWriteSyncer zapcore.WriteSyncer
		fileWriteSyncer, fileCloser, err = zap.Open(conf.Path)
		if err != nil {
			return nil, nil, err
		}
		coreLoggers = append(coreLoggers, zapcore.NewCore(encoder, fileWriteSyncer, currentLevelPriority))
	}
	if conf.Stdout {
		var stdLevelPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel && lvl >= currentLevel
		})
		stdWriteSyncer := zapcore.Lock(os.Stdout)
		coreLoggers = append(coreLoggers, zapcore.NewCore(encoder, stdWriteSyncer, stdLevelPriority))
	}
	if conf.Stderr {
		var errLevelPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel && lvl >= currentLevel
		})
		errWriteSyncer := zapcore.Lock(os.Stderr)
		coreLoggers = append(coreLoggers, zapcore.NewCore(encoder, errWriteSyncer, errLevelPriority))
	}
	core := zapcore.NewTee(coreLoggers...)
	log := zap.New(core)
	return log, fileCloser, nil
}

func newZapEncodeConfig() (ec zapcore.EncoderConfig) {
	ec.NameKey = "l"
	ec.TimeKey = "ts"
	ec.LevelKey = "lvl"
	ec.MessageKey = "msg"
	ec.LineEnding = zapcore.DefaultLineEnding
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	ec.EncodeDuration = zapcore.SecondsDurationEncoder
	ec.EncodeCaller = zapcore.ShortCallerEncoder
	return ec
}
