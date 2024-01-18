package logger

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var Log *log.Logger

func Init() error {
	now := time.Now()
	timeStr := now.Format("2006-01-02 15:04:05")
	// TODO mkdir
	filePath := "./log/" + timeStr + "-" + strconv.FormatInt(now.UnixMilli(), 10) + ".log"
	fileWriter, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	writer := io.MultiWriter(os.Stdout, fileWriter)
	prefix := "[perpetua] "
	logFlag := log.Ldate | log.Ltime | log.Lshortfile
	Log = log.New(writer, prefix, logFlag)
	return nil
}
