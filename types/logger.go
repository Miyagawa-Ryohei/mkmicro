package types

type Logger interface {
	Info(msg string)
	Infof(format string, binder ...interface{})
	Debug(msg string)
	Debugf(format string, binder ...interface{})
	Warn(msg string)
	Warnf(format string, binder ...interface{})
	Error(msg string)
	Errorf(format string, binder ...interface{})
	Flush()
}
