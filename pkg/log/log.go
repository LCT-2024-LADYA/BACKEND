package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

const ROOT = "BACKEND"

func UnitFormatter() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		path := strings.Split(file, ROOT)
		return fmt.Sprintf("%s:%d", fmt.Sprintf("%s%s", ROOT, path[len(path)-1]), line)
	}
}

func InitLoggers() (zerolog.Logger, *os.File) {
	UnitFormatter()

	loggerFile, err := os.OpenFile("log/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic("Error opening info log file")
	}

	logger := zerolog.New(loggerFile).With().Timestamp().Caller().Logger()

	return logger, loggerFile
}
