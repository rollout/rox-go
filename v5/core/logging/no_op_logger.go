package logging

type noOpLogger struct {
}

func NewNoOpLogger() Logger {
	return &noOpLogger{}
}

func (*noOpLogger) Debug(message string, err interface{}) {
}

func (*noOpLogger) Warn(message string, err interface{}) {
}

func (*noOpLogger) Error(message string, err interface{}) {
}
