package log

import (
	"io"

	charmLog "github.com/charmbracelet/log"
)

type Logger interface {
	Print(msg any)
	Printf(format string, args ...any)
	Debug(msg any)
	Debugf(format string, args ...any)
	Info(msg any)
	Infof(format string, args ...any)
	Warn(msg any)
	Warnf(format string, args ...any)
	Error(msg any)
	Errorf(format string, args ...any)
	Fatal(msg any)
	Fatalf(format string, args ...any)
}

type StandardLogger struct {
	internalLogger *charmLog.Logger
}

func New(out io.Writer, prefix string) *StandardLogger {
	internalLogger := charmLog.NewWithOptions(
		out,
		charmLog.Options{
			Prefix: prefix,
		},
	)

	return &StandardLogger{
		internalLogger,
	}
}

func (l *StandardLogger) Print(msg any) {
	l.internalLogger.Print(msg)
}

func (l *StandardLogger) Printf(format string, args ...any) {
	l.internalLogger.Printf(format, args...)
}

func (l *StandardLogger) Debug(msg any) {
	l.internalLogger.Debug(msg)
}

func (l *StandardLogger) Debugf(format string, args ...any) {
	l.internalLogger.Debugf(format, args...)
}

func (l *StandardLogger) Info(msg any) {
	l.internalLogger.Info(msg)
}

func (l *StandardLogger) Infof(format string, args ...any) {
	l.internalLogger.Infof(format, args...)
}

func (l *StandardLogger) Warn(msg any) {
	l.internalLogger.Warn(msg)
}

func (l *StandardLogger) Warnf(format string, args ...any) {
	l.internalLogger.Warnf(format, args...)
}

func (l *StandardLogger) Error(msg any) {
	l.internalLogger.Error(msg)
}

func (l *StandardLogger) Errorf(format string, args ...any) {
	l.internalLogger.Errorf(format, args...)
}

func (l *StandardLogger) Fatal(msg any) {
	l.internalLogger.Fatal(msg)
}

func (l *StandardLogger) Fatalf(format string, args ...any) {
	l.internalLogger.Fatalf(format, args...)
}
