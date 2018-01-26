package main

import "net/http"

// registrationPage - вывод статической страницы регистрации.
func registrationPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
	<html>
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
				<p>	  <button type="submit" >Регистрация</button>			</p>
				<p><font size="2" color="red" face="Arial">* - минимум 7 символов</font></p>
		  </form>
		</body>
		</html>`))
}

// authorizationPage - вывод статической страницы авторизации.
func authorizationPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
	<html>
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
				<p> <button type="submit">Вход</button>	</p>
		  </form>
		</body>
	</html>`))
}

// addNotePage - вывод статической добавления новой записи.
func addNotePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
	<html>
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
			<p><a href="/logoutHandler">Выйти из аккаунта</a> <a href="/profileHandler">Профиль</a></p>
			<form method="post" enctype="application/x-www-form-urlencoded" action="/addNoteHandler" class="center">
				<h3>Страница добавления записи</h3>
				<p>
					<label>Заголовок:</label>
					<input type="text"  name="header">
				</p>		
				<p>
					<label >Текст записи:</label>
					<textarea name="text" cols="33" rows="3"></textarea>
				</p>		
				<p>  <button type="submit" >Добавить</button>	</p>
		  </form>
		</body>
	</html>`))
}
