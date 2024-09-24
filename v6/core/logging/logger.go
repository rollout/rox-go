package logging

type Logger interface {
	Debug(message string, err interface{})
	Warn(message string, err interface{})
	Error(message string, err interface{})
}
