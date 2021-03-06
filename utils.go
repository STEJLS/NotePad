package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// toString - преобразует bson.ObjectId в строкую
func toString(s bson.ObjectId) string {
	return s.Hex()
}

// initTemplate - инициализирует html шаблоны.
func initTemplate() {
	var err error
	funcMap := template.FuncMap{
		"toString": toString,
	}

	profile, err = profile.Funcs(funcMap).ParseFiles("./templates/profile.html")
	if err != nil {
		log.Fatal("Фатал. При парсинге страницы профиля: " + err.Error())
	}

	editNote, err = editNote.Funcs(funcMap).ParseFiles("./templates/editNoteForm.html")
	if err != nil {
		log.Fatal("Фатал. При парсинге страницы редактирования записи: " + err.Error())
	}
}

// connectToDB - устанавливает соединение с БД и инициализирует глобальные переменные.
func connectToDB(host string, port int, DBName string) {
	var err error
	DBsession, err = mgo.Dial(fmt.Sprintf("mongodb://%v:%v", host, port))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Фатал. При подключении к серверу БД(%v:%v): ", host, port) + err.Error())
	}

	usersColl = DBsession.DB(DBName).C("users")
	noteColl = DBsession.DB(DBName).C("notes")

	log.Printf("Инфо. Подключение к базе данных установлено.")
}

// validatePassword - проверяется пароль на требования системы.
func validatePassword(password string) *Error {
	if len(password) < 7 {
		return &Error{http.StatusBadRequest, "Длина пароля должна быть не менее 7 символов."}
	}

	return nil
}

// checkLogin - проверяет логин на пустую строку, а так же проверяет нет ли пользователя с таким логином.
func checkLogin(login string) *Error {
	if login == "" {
		return &Error{http.StatusBadRequest, "Логин не может быть пустой строкой."}
	}

	var result User
	err := usersColl.Find(bson.M{"login": login}).One(&result)

	if err == nil {
		return &Error{http.StatusBadRequest, "Данный логин уже используется."}
	}

	if err.Error() != "not found" {
		log.Println("Ошибка. При поиске в БД пользователя(login - " + login + "): " + err.Error())
		return &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."}
	}

	return nil
}

// sendErrorPage - генерирует страницу с ошибкой и отправляет ее.
func sendErrorPage(w http.ResponseWriter, err *Error) {
	w.WriteHeader(err.Status)
	w.Write([]byte(`<!DOCTYPE html>
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>NotePad</title>				
				<style>
				.center {
					width: 400px; /* Ширина элемента в пикселах */
					padding: 10px; /* Поля вокруг текста */
					margin: auto; /* Выравниваем по центру */
					background: #fc0; /* Цвет фона */
				}
				</style>
			</head>
			<body>
				<div class="center">
					<h1>Произошла ошибка!</h1>
					<p>Статус: ` + strconv.Itoa(err.Status) + `</p>
					<p>Информация: ` + err.Text + `</p>
					<p><a href="registrationPage">Регистрация</a> <a href="authorizationPage">  Авторизация</a> <a href="profileHandler">  Профиль</a></p>
				</div>
			</body>
		</html>		
		`))
}

// generateMD5hash - хэширует пароль по правилу: md5( md5(password) + salt)
func generateMD5hash(password string) [md5.Size]byte {
	md5hash := md5.Sum([]byte(password))

	temp := make([]byte, 0, md5.Size+len(salt))
	for _, item := range md5hash {
		temp = append(temp, item)
	}
	for _, item := range salt {
		temp = append(temp, item)
	}

	return md5.Sum(temp)
}

// generateToken - генерирует уникальный токен для авторизации.
func generateToken() string {
	token, err := uuid.NewV4()

	if err != nil {
		log.Println("Ошибка. При генерации токена: " + err.Error())
	}

	return token.String()
}

// getLoginFromCookie - проверяет cookie с токеном в браузере пользователя,
// и в случае его наличия возвращает логин пользовател.
func getLoginFromCookie(w http.ResponseWriter, r *http.Request) string {
	token := getTokenFromCookie(w, r)

	if token == "" {
		return ""
	}

	lock.RLock()
	login, ok := sessions[token]
	lock.RUnlock()

	if !ok || login == "" {
		redirectPage(w, "/authorizationPage")
		return ""
	}

	return login
}

// getTokenFromCookie - возвращает token авторизации из cookie.
func getTokenFromCookie(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("token")

	if err != nil && err.Error() == "http: named cookie not present" {
		redirectPage(w, "/authorizationPage")
		return ""
	}

	if err != nil {
		log.Println("Ошибка. При чтении cookie: " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return ""
	}

	if cookie.Value == "" {
		redirectPage(w, "/authorizationPage")
		return ""
	}

	return cookie.Value
}

// redirectPage - перенаправляет на другую страницу.
func redirectPage(w http.ResponseWriter, url string) {
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusFound)
}
