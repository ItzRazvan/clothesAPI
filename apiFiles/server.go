package apiFiles

import (
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

// porneste serverului cu toate functiile sale
func ServerStart() {
	app := echo.New()

	//initializare template randerer
	t := &Template{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}
	app.Renderer = t
	app.GET("/", renderIndex)
	app.GET("/genKey", generateKey)

	app.GET("/login", renderLogin)
	app.GET("/signin", renderSingin)

	app.POST("/login", loginTry)
	app.POST("/signin", signinTry)

	app.Logger.Fatal(app.Start(":8080"))
}

// Afiseaza index.html cand este accesata pagina principala
func renderIndex(c echo.Context) error {

	//cheie := existaSauGenereazaCheie()

	key := map[string]interface{}{
		"Cheie": "1234567890",
	}

	if isLoggedIn(c) {
		//randeaza pagina principala\
		return c.Render(http.StatusOK, "index.html", key)
	} else {
		//randeaza pagina de login
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return nil
	}
}

// Functie care randeaza pagina de login
func renderLogin(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.Render(http.StatusOK, "login.html", nil)
	}
	http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)
	return nil
}

// Functie care randeaza pagina de signin
func renderSingin(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.Render(http.StatusOK, "signin.html", nil)
	}
	http.Redirect(c.Response(), c.Request(), "/", http.StatusSeeOther)
	return nil
}

// Genereaza o cheie noua dupa ce butonul de generare este apasat
func generateKey(c echo.Context) error {
	schimbareCheie()
	return nil
}
