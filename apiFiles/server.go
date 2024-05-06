package apiFiles

import (
	"encoding/json"
	"io"
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
	app.GET("/genKey", getKey)
	app.GET("/logout", logout)

	app.GET("/login", renderLogin)
	app.GET("/signin", renderSingin)

	app.POST("/login", loginTry)
	app.POST("/signin", signinTry)

	app.GET("/apiTest", renderApiTest)

	//Functiile din api.go
	app.GET("/api/haine", returneazaHaine)
	app.POST("/api/POST", posteazaHaine)
	app.GET("/api/haine:id", returneazaHainaDupaId)

	app.Logger.Fatal(app.Start(":8080"))
}

// Afiseaza index.html cand este accesata pagina principala
func renderIndex(c echo.Context) error {
	//verifica daca userul este logat
	if isLoggedIn(c) {
		cheie := getCheieFromDB(c)

		key := map[string]interface{}{
			"Cheie": cheie,
		}

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
/*func generateKey(c echo.Context) error {
	schimbareCheie(c, 0)
	refreshPage(c)
	return nil
}*/

// funcite care genereaza cheie noua
func getKey(c echo.Context) error {
	schimbareCheie(c, 0)
	cheie := getCheieFromDB(c)
	js, err := json.Marshal(cheie)
	check(err)
	c.Response().Writer.Header().Set("Content-Type", "application/json")
	c.Response().Writer.Write(js)
	return nil
}

// Functie pentru apasarea butonului de logout
func logout(c echo.Context) error {
	//stergem cookie ul
	removeCookie(c)

	//redirectam catre pagina de login
	return c.Redirect(http.StatusSeeOther, "/login")
}

// Functie care randeaza pagina de testare a api ului
func renderApiTest(c echo.Context) error {
	//vom trimite toate hainele in pagina

	//incepem prin a face un rqeuest catre api pentru a lua toate hainele
	cheie := getCheieFromDB(c)
	if cheie == "" {
		c.Redirect(302, "/login")
		return nil
	}
	linkPtReq := "http://localhost:8080/api/haine?key=" + cheie
	resp, err := http.Get(string(linkPtReq))
	check(err)
	body, err := io.ReadAll(resp.Body)
	check(err)
	resp.Body.Close()

	// Unmarshal the JSON response
	var haine []haina
	err = json.Unmarshal(body, &haine)
	check(err)

	// Now you can use 'haine' in your template
	return c.Render(http.StatusOK, "apiTest.html", map[string]interface{}{
		"haine": haine,
	})
}
