package logger

import (
	"github.com/IUnlimit/perpetua/internal/conf"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

func Init() {
	initLog()
}

func initLog() {
	config := conf.Config
	rotateOptions := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour * 24),
	}
	rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(config.Log.LogAging))
	if config.Log.LogForceNew {
		rotateOptions = append(rotateOptions, rotatelogs.ForceNewFile())
	}

	w, err := rotatelogs.New(path.Join(config.ParentPath+"/logs", "%Y-%m-%d.log"), rotateOptions...)
	if err != nil {
		log.Errorf("rotatelogs init err: %v", err)
		panic(err)
	}

	consoleFormatter := LogFormat{EnableColor: config.Log.LogColorful}
	fileFormatter := LogFormat{EnableColor: false}
	Hook = NewLocalHook(w, consoleFormatter, fileFormatter, GetLogLevel(config.Log.LogLevel)...)
	log.AddHook(Hook)
}
