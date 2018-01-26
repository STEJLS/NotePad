package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/STEJLS/notepad/XMLconfig"
)

func main() {
	InitFlags()
	initTemplate()
	logFile := InitLogger(logSource)
	defer logFile.Close()

	config := XMLconfig.Parse(configSource)

	connectToDB(config.Db.Host, config.Db.Port, config.Db.Name)
	defer DBsession.Close()

	server := http.Server{
		Addr: fmt.Sprintf("%v:%v", config.HTTP.Host, config.HTTP.Port),
	}

	http.HandleFunc("/registrationPage", registrationPage)
	http.HandleFunc("/authorizationPage", authorizationPage)
	http.HandleFunc("/addNotePage", addNotePage)
	http.HandleFunc("/addNoteHandler", addNoteHandler)
	http.HandleFunc("/registrationHandler", registrationHandler)
	http.HandleFunc("/authorizationHandler", authorizationHandler)
	http.HandleFunc("/logoutHandler", logoutHandler)
	http.HandleFunc("/profileHandler", profileHandler)
	http.HandleFunc("/deleteNoteHandler", deleteNoteHandler)
	http.HandleFunc("/editNoteFormHandler", editNoteFormHandler)
	http.HandleFunc("/saveEditingHandler", saveEditingHandler)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}
