package job

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

var (
	path, _    = os.Getwd()
	moduleName = "opensnail.com/snail-job"
)

type DefaultFormatter struct {
	ForceColors bool
}

func trimProjectPath(fullPath, projectRoot string) string {
	relativePath, err := filepath.Rel(projectRoot, fullPath)
	if err != nil {
		// 如果出错，直接返回原路径
		return fullPath
	}
	return relativePath
}

func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	d := entry.Data

	if entry.HasCaller() {
		// 自定义格式：时间、级别、调用者、消息
		log := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s:%d] - %s\n",
			entry.Time.Format("2006-01-02 15:04:05.000"),
			d["logger"],
			strings.ToUpper(entry.Level.String()),
			//filepath.Base(path)+"/"+trimProjectPath(entry.Caller.File, path),
			trimProjectPath(entry.Caller.File, path),
			trimProjectPath(entry.Caller.Function, moduleName),
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
