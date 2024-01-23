package logger

import (
	global "github.com/IUnlimit/perpetua/internal"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

func Init() {
	initLog()
}

func initLog() {
	config := global.Config
	rotateOptions := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour * 24),
	}
	rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(config.Log.Aging))
	if config.Log.ForceNew {
		rotateOptions = append(rotateOptions, rotatelogs.ForceNewFile())
	}

	w, err := rotatelogs.New(path.Join(config.ParentPath+"/logs", "%Y-%m-%d.log"), rotateOptions...)
	if err != nil {
		log.Errorf("rotatelogs init err: %v", err)
		panic(err)
	}

	levels := GetLogLevel(config.Log.Level)
	log.SetLevel(levels[0]) // hook levels doesn't work
	consoleFormatter := LogFormat{EnableColor: config.Log.Colorful}
	fileFormatter := LogFormat{EnableColor: false}
	Hook = NewLocalHook(w, consoleFormatter, fileFormatter, levels...)
	log.AddHook(Hook)
}
