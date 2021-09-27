package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Log struct {
	msg     string
	traceId string
	logger  *Logger
}

var logPool = sync.Pool{
	New: func() interface{} {
		return new(Log)
	},
}

func getLog() *Log {
	return logPool.Get().(*Log)
}

func putLog(log *Log) {
	log.reset()
	logPool.Put(log)
}

func (l *Log) reset() {
	l.msg = ""
	l.traceId = ""
	l.logger = nil
}

func (l *Log) TraceId(reqId string) *Log {
	if l == nil {
		return nil
	}
	l.traceId = reqId
	return l
}

func (l *Log) AddMsg(msg string) *Log {
	if l == nil {
		return nil
	}
	if len(l.msg) == 0 {
		l.msg = msg
	} else {
		l.msg = l.msg + ": " + msg
	}
	return l
}

func (l *Log) AddMsgf(format string, params ...interface{}) *Log {
	return l.AddMsg(fmt.Sprintf(format, params...))
}

func (l *Log) Info(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.Info(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.Info(l.msg)
	}
	putLog(l)
}

func (l *Log) Infof(format string, params ...interface{}) {
	l.Info(fmt.Sprintf(format, params...))
}

func (l *Log) Error(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.Error(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.Error(l.msg)
	}
	putLog(l)
}

func (l *Log) Errorf(format string, params ...interface{}) {
	l.Error(fmt.Sprintf(format, params...))
}

func (l *Log) Debug(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.Debug(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.Debug(l.msg)
	}
	putLog(l)
}

func (l *Log) Debugf(format string, params ...interface{}) {
	l.Debug(fmt.Sprintf(format, params...))
}

func (l *Log) Warning(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.Warn(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.Warn(l.msg)
	}
	putLog(l)
}

func (l *Log) Warningf(format string, params ...interface{}) {
	l.Warning(fmt.Sprintf(format, params...))
}

func (l *Log) Panic(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.Panic(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.Panic(l.msg)
	}
	putLog(l)
}

func (l *Log) Panicf(format string, params ...interface{}) {
	l.Panic(fmt.Sprintf(format, params...))
}

func (l *Log) DPanic(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.DPanic(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.DPanic(l.msg)
	}
	putLog(l)
}

func (l *Log) DPanicf(format string, params ...interface{}) {
	l.DPanic(fmt.Sprintf(format, params...))
}

func (l *Log) Fatal(msg string) {
	if l == nil {
		return
	}
	if l.logger.isLock {
		l.logger.mux.Lock()
		defer l.logger.mux.Unlock()
	}
	l.AddMsg(msg)
	if len(l.traceId) > 0 {
		l.logger.zapLogger.Fatal(l.msg, zap.String(traceIdKey, l.traceId))
	} else {
		l.logger.zapLogger.Fatal(l.msg)
	}
	putLog(l)
}

func (l *Log) Fatalf(format string, params ...interface{}) {
	l.Fatal(fmt.Sprintf(format, params...))
}
