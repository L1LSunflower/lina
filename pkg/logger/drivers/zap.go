package drivers

import (
	"context"

	"github.com/L1LSunflower/lina/pkg/logger"
)

type ZapLogger struct {
}

func (z ZapLogger) Debug(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Debugf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) DebugWithCtx(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Info(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Infof(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) InfoWithCtx(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Warn(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Warnf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) WarnWithCtx(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Error(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Errorf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) ErrorWithCtx(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Fatal(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Fatalf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) FatalWithCtx(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Panic(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) Panicf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z ZapLogger) PanicWithCtx(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func NewZapLogger() logger.Logger {
	return new(ZapLogger)
}
