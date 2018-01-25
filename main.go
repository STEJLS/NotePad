package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/STEJLS/notepad/XMLconfig"
)

func main() {
	InitFlags()
	logFile := InitLogger(logSource)
	defer logFile.Close()

	config := XMLconfig.Parse(configSource)

	connectToDB(config.Db.Host, config.Db.Port, config.Db.Name)
	defer DBsession.Close()

	server := http.Server{
		Addr: fmt.Sprintf("%v:%v", config.HTTP.Host, config.HTTP.Port),
	}

	http.HandleFunc("/registrationPage", registrationPage)
	http.HandleFunc("/registrationHandler", registrationHandler)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}
