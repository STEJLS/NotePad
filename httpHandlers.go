package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// registrationHandler - обработчик, который осуществляет регистрацию нового пользователя.
// Принимет Post запрос с переменными "login" и "password".
func registrationHandler(w http.ResponseWriter, r *http.Request) {
	login := strings.ToLower(r.FormValue("login"))
	password := r.FormValue("password")

	if err := validatePassword(password); err != nil {
		sendErrorPage(w, err)
		return
	}

	if err := checkLogin(login); err != nil {
		sendErrorPage(w, err)
		return
	}

	err := usersColl.Insert(&User{Login: login, Password: generateMD5hash(password)})
	if err != nil {
		log.Println("Ошибка. При вставке в БД пользователя: " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	log.Println("Инфо. Пользователь " + login + " зарегистрировался.")

	redirectPage(w, "/authorizationPage")
}

// authorizationHandler - обработчик, который осуществляет авторизация пользователя. Токен записывается в cookie.
func authorizationHandler(w http.ResponseWriter, r *http.Request) {
	login := strings.ToLower(r.FormValue("login"))
	password := r.FormValue("password")

	if err := validatePassword(password); err != nil {
		sendErrorPage(w, err)
		return
	}

	if login == "" {
		sendErrorPage(w, &Error{http.StatusBadRequest, "Логин не может быть пустой строкой."})
		return
	}

	var user User
	err := usersColl.Find(bson.M{"login": login}).One(&user)
	if err != nil && err.Error() == "not found" {
		sendErrorPage(w, &Error{http.StatusBadRequest, "Пользователя с таким логином не существует."})
		return
	}

	if err != nil {
		log.Println("Ошибка. При поиске в БД пользователя(логин - " + login + " ): " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	estimatePass := generateMD5hash(password)

	if !bytes.Equal(user.Password[:], estimatePass[:]) {
		sendErrorPage(w, &Error{http.StatusBadRequest, "Неверный пароль."})
		return
	}

	token := generateToken()
	sessions[token] = login

	http.SetCookie(w, &http.Cookie{Name: "token", Value: token})

	log.Println("Инфо. Пользователь " + login + " авторизовался.")
	redirectPage(w, "/profileHandler")
}

// addNoteHandler - добавляет новую заметку.
func addNoteHandler(w http.ResponseWriter, r *http.Request) {
	login := getLoginFromCookie(w, r)

	if login == "" {
		return
	}

	note := Note{
		Author:      login,
		Header:      r.FormValue("header"),
		Text:        r.FormValue("text"),
		CreatedDate: time.Now(),
	}

	if err := noteColl.Insert(&note); err != nil {
		log.Println("Ошибка. При вставке в БД новой заметки пользователя(логин - " + login + " ): " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	redirectPage(w, "/profileHandler")
}

// logoutHandler - реализует выход из аккаунта
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromCookie(w, r)

	if token == "" {
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "token", Expires: time.Now().UTC()})
	delete(sessions, token)

	redirectPage(w, "/authorizationPage")
}

// profileHandler - генерирует страницу профиля пользователя с его личными записями
func profileHandler(w http.ResponseWriter, r *http.Request) {
	login := getLoginFromCookie(w, r)

	if login == "" {
		return
	}

	var result []Note

	err := noteColl.Find(bson.M{"author": login}).Sort("-createdDate").All(&result)

	if err != nil {
		log.Println("Ошибка. При поиске в БД заметок пользователя(логин - " + login + " ): " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	profile.Execute(w, result)
}

// deleteNoteHandler - удаляет заметку пользователя.
func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	login := getLoginFromCookie(w, r)

	if login == "" {
		return
	}

	id := r.FormValue("id")
	if !bson.IsObjectIdHex(id) {
		log.Println("Ошибка. При удалении заметки некорректный id =  " + id)
		sendErrorPage(w, &Error{http.StatusBadRequest, "Получен некорректный id."})
		return
	}

	err := noteColl.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Println("Ошибка. При удалении заметки ( id =  " + id + " )в БД: " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	redirectPage(w, "/profileHandler")

	log.Println("Инфо. Заметка с id = '" + id + "' удалена.")
}

// editNoteFormHandler - генерирует страницу для редактирования заметки.
func editNoteFormHandler(w http.ResponseWriter, r *http.Request) {
	login := getLoginFromCookie(w, r)

	if login == "" {
		return
	}

	id := r.FormValue("id")
	if !bson.IsObjectIdHex(id) {
		log.Println("Ошибка. При редактировании заметки некорректный id =  " + id)
		sendErrorPage(w, &Error{http.StatusBadRequest, "Получен некорректный id."})
		return
	}

	var result Note

	err := noteColl.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil && err.Error() != "not found" {
		log.Println("Ошибка. При поиске в БД заметки с id = " + id + " " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	if err != nil {
		sendErrorPage(w, &Error{http.StatusBadRequest, "С указанным id нет заметки."})
		return
	}

	editNote.Execute(w, result)
}

// saveEditingHandler - осуществляет изменение заметки.
func saveEditingHandler(w http.ResponseWriter, r *http.Request) {
	login := getLoginFromCookie(w, r)

	if login == "" {
		return
	}

	id := r.FormValue("id")
	if !bson.IsObjectIdHex(id) {
		log.Println("Ошибка. При сохранении редактирования некорректный id =  " + id)
		sendErrorPage(w, &Error{http.StatusBadRequest, "Получен некорректный id."})
		return
	}

	err := noteColl.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"header": r.FormValue("header"), "text": r.FormValue("text")}})
	if err != nil {
		log.Println("Ошибка. При редактировании в БД заметки с id = " + id + " " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	redirectPage(w, "/profileHandler")

	log.Println("Инфо. Заметка с id = '" + id + "' отредактирована.")
}
