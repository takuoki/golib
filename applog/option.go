package applog

// Option is an option for logger generation.
type Option func(Logger) error

// LevelOption sets the log level that the logger outputs.
func LevelOption(lv Level) Option {
	return func(l Logger) error {
		return l.setLevel(lv)
	}
}

// TimeFormatOption sets the time format that the logger outputs.
func TimeFormatOption(format string) Option {
	return func(l Logger) error {
		return l.setTimeFormat(format)
	}
}

// ImageTagOption sets the image tag that the logger outputs.
func ImageTagOption(tag string) Option {
	return func(l Logger) error {
		return l.setImageTag(tag)
	}
}
