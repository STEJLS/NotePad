package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
)

// InitFlags - инициализирует флаги командной строки.
func InitFlags() {
	flag.StringVar(&logSource, "log_source", "log.txt", "Source for log file")
	flag.StringVar(&configSource, "config_source", "config.xml", "Source for config file")
	flag.Parse()
}

// InitLogger - инициализирует логгер.
func InitLogger(destination string) *os.File {
	logfile, err := os.OpenFile(logSource, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("Фатал. Файл логов (%q) не открылся: ", logSource, err)
	}

	log.SetOutput(logfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return logfile
}

// connectToDB - устанавливает соединение с БД и инициализирует глобальные переменные.
func connectToDB(host string, port int, DBName string) {
	var err error
	noteDBsession, err = mgo.Dial(fmt.Sprintf("mongodb://%v:%v", host, port))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Фатал. При подключении к серверу БД(%v:%v): ", host, port) + err.Error())
	}

	log.Printf("Инфо. Подключение к базе данных установлено.")
}
