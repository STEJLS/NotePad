package main

import (
	"log"
	"net/http"
	"strings"
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
			<form method="post" enctype="application/x-www-form-urlencoded" action="/registrationHandler" class="center">
			<p>
			  <label>Логин:</label>
			  <input type="text"  name="login">
			</p>		
			<p>
			  <p><font size="2" color="red" face="Arial">минимум 7 символов</font></p>
			  <label >Пароль:</label>
			  <input type="password"  name="password">
			</p>		
			<p class="login-submit">
			  <button type="submit" class="login-button">Регистрация</button>
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
