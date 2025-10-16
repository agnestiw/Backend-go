package config

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	})
	log.Println("File logger berhasil diinisialisasi.")
}