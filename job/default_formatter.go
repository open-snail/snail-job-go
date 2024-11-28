package job

import (
	"fmt"
	"github.com/open-snail/snail-job-go/util"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	path, _    = os.Getwd()
	moduleName = "github.com/open-snail"
)

type DefaultFormatter struct {
	ForceColors bool
}

func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	d := entry.Data

	if entry.HasCaller() {
		// 自定义格式：时间、级别、调用者、消息
		log := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s:%d] - %s\n",
			entry.Time.Format("2006-01-02 15:04:05.000"),
			d["logger"],
			strings.ToUpper(entry.Level.String()),
			util.TrimProjectPath(entry.Caller.File, path),
			util.TrimProjectPath(entry.Caller.Function, moduleName),
			entry.Caller.Line,
			entry.Message,
		)
		return []byte(log), nil
	} else {
		// 自定义格式：时间、级别、调用者、消息
		log := fmt.Sprintf("[%s] [%s] [%s] \n%s\n",
			entry.Time.Format("2006-01-02 15:04:05.000"),
			d["logger"],
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		)
		return []byte(log), nil
	}

}
