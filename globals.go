package main

import mgo "gopkg.in/mgo.v2"

// logFileName - имя файла для логов, задается через флаг командной строки.
var logSource string

// ConfigSource - имя файла для конфига, задается через флаг командной строки.
var configSource string

// noteDBsession - указатель на сессию подключения к БД notepad.
var noteDBsession *mgo.Session
