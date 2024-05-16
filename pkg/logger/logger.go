package logger

import "context"

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugWithCtx(ctx context.Context, format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoWithCtx(ctx context.Context, format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnWithCtx(ctx context.Context, format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorWithCtx(ctx context.Context, format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	FatalWithCtx(ctx context.Context, format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	PanicWithCtx(ctx context.Context, format string, args ...interface{})
}
