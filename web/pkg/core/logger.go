package core

type Logger interface {
	Fatal(f string, args ...interface{})
	Error(f string, args ...interface{})
	Warn(f string, args ...interface{})
	Info(f string, args ...interface{})
	Debug(f string, args ...interface{})
}
