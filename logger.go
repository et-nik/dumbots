package main

import (
	metamod "github.com/et-nik/metamod-go"
)

type Logger struct {
	MetaUtilFuncs *metamod.MUtilFuncs
}

func NewLogger(metaUtilFuncs *metamod.MUtilFuncs) *Logger {
	return &Logger{
		MetaUtilFuncs: metaUtilFuncs,
	}
}

func (l *Logger) Message(message string) {
	if l.MetaUtilFuncs == nil {
		return
	}

	l.MetaUtilFuncs.LogMessage(message)
}

func (l *Logger) Messagef(format string, args ...interface{}) {
	if l.MetaUtilFuncs == nil {
		return
	}

	l.MetaUtilFuncs.LogMessagef(format, args...)
}

func (l *Logger) Error(message string) {
	if l.MetaUtilFuncs == nil {
		return
	}

	l.MetaUtilFuncs.LogError(message)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.MetaUtilFuncs == nil {
		return
	}

	l.MetaUtilFuncs.LogErrorf(format, args...)
}

func (l *Logger) Debug(message string) {
	if l.MetaUtilFuncs == nil {
		return
	}

	l.MetaUtilFuncs.LogDeveloper(message)
}

func (l *Logger) Debugf(message string, args ...interface{}) {
	if l.MetaUtilFuncs == nil {
		return
	}

	l.MetaUtilFuncs.LogDeveloperf(message, args...)
}
