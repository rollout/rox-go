package logging

var (
	mainLogger  Logger
	dummyLogger = NewNoOpLogger()
)

func SetLogger(logger Logger) {
	mainLogger = logger
}

func GetLogger() Logger {
	if mainLogger != nil {
		return mainLogger
	}

	return dummyLogger
}
