package applog

// Level is the log level.
// Higher levels are more important.
type Level int

// Level list. Default is InfoLevel.
const (
	UnknownLevel Level = iota - 3
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

// ParseLevel returns the log level based on the string.
func ParseLevel(lv string) Level {
	switch lv {
	case "critical", "CRITICAL", "Critical":
		return CriticalLevel
	case "error", "ERROR", "Error":
		return ErrorLevel
	case "warn", "WARN", "Warn":
		return WarnLevel
	case "info", "INFO", "Info":
		return InfoLevel
	case "debug", "DEBUG", "Debug":
		return DebugLevel
	case "trace", "TRACE", "Trace":
		return TraceLevel
	}
	return UnknownLevel
}

// String returns a string of the log level.
func (lv Level) String() string {
	switch lv {
	case CriticalLevel:
		return "CRITICAL"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

func shouldPrint(setting, print Level) bool {
	return setting <= print
}
