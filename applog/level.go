package applog

import "errors"

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
// If an undefined string is specified, an error will be returned.
func ParseLevel(lv string) (Level, error) {
	switch lv {
	case "critical", "CRITICAL", "Critical":
		return CriticalLevel, nil
	case "error", "ERROR", "Error":
		return ErrorLevel, nil
	case "warn", "WARN", "Warn":
		return WarnLevel, nil
	case "info", "INFO", "Info":
		return InfoLevel, nil
	case "debug", "DEBUG", "Debug":
		return DebugLevel, nil
	case "trace", "TRACE", "Trace":
		return TraceLevel, nil
	}
	return UnknownLevel, errors.New("invalid string for the log level")
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
