package service

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
)

var Logger = log.Init(&log.Config{Level: log.TraceLevel, Format: log.FormatterJson})

func LF() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(Logger).Srv("trading").Node("trading")
	}
}

func L() log.CLogger {
	return LF()()
}
