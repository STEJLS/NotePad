package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// registrationPage - вывод статической страницы регистрации.
func registrationPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
		<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
			  <title>NotePad</title>
			  <style>
				  .center {
				   width: 250px; /* Ширина элемента в пикселах */
				   padding: 10px; /* Поля вокруг текста */
				   margin: auto; /* Выравниваем по центру */
				   background: #fc0; /* Цвет фона */
				  }
				 </style>
		</head>
		<body>
		<p><a href="/authorizationPage">Авторизация</a></p>
			<form method="post" enctype="application/x-www-form-urlencoded" action="/registrationHandler" class="center">
			<h3>Страница регистрации</h3>
			<p>
			  <label>Логин:</label>
			  <input type="text"  name="login">
			</p>		
			<p>

			  <label >Пароль*:</label>
			  <input type="password"  name="password">
			</p>		
			<p class="login-submit">
			  <button type="submit" class="login-button">Регистрация</button>
			</p>
			<p><font size="2" color="red" face="Arial">* - минимум 7 символов</font></p>

		  </form>
		</body>
		</html>		
		`))
}

// authorizationPage - вывод статической страницы авторизации.
func authorizationPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
		<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
			  <title>NotePad</title>
			  <style>
				  .center {
				   width: 250px; /* Ширина элемента в пикселах */
				   padding: 10px; /* Поля вокруг текста */
				   margin: auto; /* Выравниваем по центру */
				   background: #fc0; /* Цвет фона */
				  }
				 </style>
		</head>
		<body>
		<p><a href="/registrationPage">Регистрация</a></p>
			<form method="post" enctype="application/x-www-form-urlencoded" action="/authorizationHandler" class="center">
			<h3>Страница авторизации</h3>
			<p>
			  <label>Логин:</label>
			  <input type="text"  name="login">
			</p>		
			<p>
			  <label >Пароль:</label>
			  <input type="password"  name="password">
			</p>		
			<p class="login-submit">
			  <button type="submit" class="login-button">Вход</button>
			</p>
		  </form>
		</body>
		</html>		
		`))
}

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

	log.Println("Инфо. Добавлен новый пользователь - " + login)
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
	answer := bytes.Equal(user.Password[:], estimatePass[:])

	if !answer {
		sendErrorPage(w, &Error{http.StatusBadRequest, "Неверный пароль."})
		return
	}

	token := generateToken()
	sessions[token] = login

	http.SetCookie(w, &http.Cookie{Name: "token", Value: token})

	log.Println("Инфо. Пользователь " + login + " авторизовался.")
}

func test(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err != nil && err.Error() == "http: named cookie not present" {
		http.Redirect(w, r, "/authorizationPage", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Println("Ошибка. При чтении cookie: " + err.Error())
		sendErrorPage(w, &Error{http.StatusInternalServerError, "Неполадки на сервере, повторите попытку позже."})
		return
	}

	if cookie.Value == "" {
		http.Redirect(w, r, "/authorizationPage", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("твой токен - " + cookie.Value))
	w.Write([]byte("\nтвой логин - " + sessions[cookie.Value]))

}
