package logger

import "fmt"

type Logger struct {
	Type  string `yaml:"Type"`
	Level string `yaml:"Level"`
}

const (
	levelDebugValue = iota
	levelInfoValue
	levelWarnValue
	levelErrorValue
	levelFatalfValue
)

const (
	levelDebug  = "debug"
	levelInfo   = "info"
	levelWarn   = "warn"
	levelError  = "error"
	levelFatalf = "fatal"
)

var levelMap = map[string]int{
	levelDebug:  levelDebugValue,
	levelInfo:   levelInfoValue,
	levelWarn:   levelWarnValue,
	levelError:  levelErrorValue,
	levelFatalf: levelFatalfValue,
}

func (l *Logger) Debug(f string, args ...interface{}) {
	if levelMap[l.Level] == levelDebugValue {
		fmt.Println()
	}
}
