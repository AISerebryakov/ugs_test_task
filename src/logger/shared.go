package logger

import (
	"sync"
)

var (
	sharedLogger *Logger
	once         sync.Once
)

func Init(conf Config) (err error) {
	once.Do(func() {
		sharedLogger, err = New(conf)
	})
	return err
}

func Close() {
	sharedLogger.Close()
}

func ReopenFile() (err error) {
	return sharedLogger.ReopenFile()
}

func TraceId(id string) *Log {
	return sharedLogger.TraceId(id)
}

func Msg(msg string) *Log {
	return sharedLogger.Msg(msg)
}

func Info(msg string) {
	sharedLogger.Info(msg)
}

func Infof(format string, params ...interface{}) {
	sharedLogger.Infof(format, params...)
}

func Error(msg string) {
	sharedLogger.Error(msg)
}

func Errorf(format string, params ...interface{}) {
	sharedLogger.Errorf(format, params...)
}

func Debug(msg string) {
	sharedLogger.Debug(msg)
}

func Debugf(format string, params ...interface{}) {
	sharedLogger.Debugf(format, params...)
}

func Warning(msg string) {
	sharedLogger.Warning(msg)
}

func Warningf(format string, params ...interface{}) {
	sharedLogger.Warningf(format, params...)
}

func Panic(msg string) {
	sharedLogger.Panic(msg)
}

func Panicf(format string, params ...interface{}) {
	sharedLogger.Panicf(format, params...)
}

func DPanic(msg string) {
	sharedLogger.DPanic(msg)
}

func DPanicf(format string, params ...interface{}) {
	sharedLogger.DPanicf(format, params...)
}

func Fatal(msg string) {
	sharedLogger.Fatal(msg)
}

func Fatalf(format string, params ...interface{}) {
	sharedLogger.Fatalf(format, params...)
}
