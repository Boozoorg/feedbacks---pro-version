package logging

import (
	"feedbacks/models"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

// InitLogger - запускает логирование в файле logs.txt/ starts logging in the logs.txt file
func InitLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   models.Conf.Logger.Logfile, // путь и файл
		MaxSize:    models.Conf.Logger.MaxSize, // megabytes
		MaxBackups: models.Conf.Logger.MaxBackups,
		MaxAge:     models.Conf.Logger.MaxAge, //days
		Compress:   true,                          // disabled by default
	})
	log.Println( fmt.Sprintf("logging into directory: %s", models.Conf.Logger.Logfile))
}
