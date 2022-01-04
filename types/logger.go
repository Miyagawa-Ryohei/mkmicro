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


type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
