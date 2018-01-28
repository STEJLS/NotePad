package main

import (
	"html/template"
	"sync"

	mgo "gopkg.in/mgo.v2"
)

// logFileName - имя файла для логов, задается через флаг командной строки.
var logSource string

// ConfigSource - имя файла для конфига, задается через флаг командной строки.
var configSource string

// DBsession - указатель на сессию подключения к серверу БД.
var DBsession *mgo.Session

// userColl - указатель на коллекции Users базы данных notepad.
var usersColl *mgo.Collection

// noteColl - указатель на коллекции Notes базы данных notepad.
var noteColl *mgo.Collection

// salt - соль для пароля.
var salt = [12]byte{152, 123, 2, 1, 6, 84, 216, 35, 140, 158, 69, 128}

// sessions - карта для авторизации пользователей. Ключ токен, а значение - логин.
var sessions = make(map[string]string)

// lock - Мьютекс для корректной параллельной работы с картой sessions.
var lock = new(sync.RWMutex)

// profile - шаблон для отображения страницы профиля.
var profile = template.New("profile.html")

// editNote - шаблон для отображения страницы редактирования заметки.
var editNote = template.New("editNoteForm.html")
