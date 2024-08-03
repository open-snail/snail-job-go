package job

import (
	"github.com/sirupsen/logrus"
	"time"
)

var (
	LocalLog = logrus.New()
)

type LogRecord struct {
	TimeStamp time.Time
	Level     string
	Thread    string
	Message   string
	Module    string
	FuncName  string
	Lineno    int
	ExcInfo   error
}
