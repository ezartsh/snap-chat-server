package logger

var AppLog LogHandler

func InitAppLog() {
	AppLog = newAppLogger()
}

func newAppLogger() LogHandler {
	return NewConsole()
}
