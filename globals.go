package main

import mgo "gopkg.in/mgo.v2"

// logFileName - имя файла для логов, задается через флаг командной строки.
var logSource string

// ConfigSource - имя файла для конфига, задается через флаг командной строки.
var configSource string

// DBsession - указатель на сессию подключения к серверу БД.
var DBsession *mgo.Session

// userColl - указатель на коллекции Users базы данных notepad.
var usersColl *mgo.Collection

// salt - соль для пароля.
var salt = [12]byte{152, 123, 2, 1, 6, 84, 216, 35, 140, 158, 69, 128}

// sessions - карта для авторизации пользователей. Ключ токен, а значение - логин.
var sessions = make(map[string]string)
