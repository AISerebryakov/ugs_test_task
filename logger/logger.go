package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

const (
	traceIdKey = "trace_id"
)

type Config struct {
	Path   string
	Stdout bool
	Stderr bool
	Level  Level
}

type Logger struct {
	zapLogger  *zap.Logger
	fileCloser func()
	config     Config
	isLock     bool
	mux        sync.Mutex
}

func New(conf Config) (_ *Logger, err error) {
	l := new(Logger)
	l.config = conf
	l.zapLogger, l.fileCloser, err = newZapFileLogger(conf)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Logger) ReopenFile() (err error) {
	if l == nil || len(l.config.Path) == 0 {
		return nil
	}
	l.isLock = true
	l.mux.Lock()
	defer l.mux.Unlock()
	l.close()
	l.zapLogger, l.fileCloser, err = newZapFileLogger(l.config)
	if err != nil {
		l.isLock = false
		return err
	}
	l.isLock = false
	return nil
}

func (l *Logger) Close() {
	if l == nil {
		return
	}
	l.isLock = true
	l.mux.Lock()
	l.close()
	l.mux.Unlock()
	l.isLock = false
}

func (l *Logger) close() {
	if l.zapLogger != nil {
		l.zapLogger.Sync()
	}
	if l.fileCloser != nil {
		l.fileCloser()
	}
	l.zapLogger = zap.NewNop()
}

//Deprecated
func (l *Logger) ReqId(id string) *Log {
	if l == nil {
		return nil
	}
	log := getLog()
	log.logger = l
	log.traceId = id
	return log
}

func (l *Logger) TraceId(id string) *Log {
	if l == nil {
		return nil
	}
	log := getLog()
	log.logger = l
	log.traceId = id
	return log
}

func (l *Logger) Msg(msg string) *Log {
	if l == nil {
		return nil
	}
	log := getLog()
	log.logger = l
	log.msg = msg
	return log
}

func (l *Logger) Msgf(format string, params ...interface{}) *Log {
	return l.Msg(fmt.Sprintf(format, params...))
}

func (l *Logger) Info(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.Info(msg)
}

func (l *Logger) Infof(format string, params ...interface{}) {
	l.Info(fmt.Sprintf(format, params...))
}

func (l *Logger) Error(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.Error(msg)
}

func (l *Logger) Errorf(format string, params ...interface{}) {
	l.Error(fmt.Sprintf(format, params...))
}

func (l *Logger) Debug(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.Debug(msg)
}

func (l *Logger) Debugf(format string, params ...interface{}) {
	l.Debug(fmt.Sprintf(format, params...))
}

func (l *Logger) Warning(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.Warn(msg)
}

func (l *Logger) Warningf(format string, params ...interface{}) {
	l.Warning(fmt.Sprintf(format, params...))
}

func (l *Logger) Panic(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.Panic(msg)
}

func (l *Logger) Panicf(format string, params ...interface{}) {
	l.Panic(fmt.Sprintf(format, params...))
}

func (l *Logger) DPanic(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.DPanic(msg)
}

func (l *Logger) DPanicf(format string, params ...interface{}) {
	l.DPanic(fmt.Sprintf(format, params...))
}

func (l *Logger) Fatal(msg string) {
	if l == nil {
		return
	}
	if l.isLock {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	l.zapLogger.Fatal(msg)
}

func (l *Logger) Fatalf(format string, params ...interface{}) {
	l.Fatal(fmt.Sprintf(format, params...))
}
