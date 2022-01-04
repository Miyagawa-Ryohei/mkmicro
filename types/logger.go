package types

type Logger interface {
	Info(format string, binder ...interface{})
	Debug(format string, binder ...interface{})
	Warn(format string, binder ...interface{})
	Error(format string, binder ...interface{})
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
