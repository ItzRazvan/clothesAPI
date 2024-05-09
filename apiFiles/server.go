package apiFiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// porneste serverului cu toate functiile sale
func ServerStart() {
	app := echo.New()

	//initializare template randerer
	t := &Template{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
	}))

	app.Renderer = t
	app.GET("/", renderIndex)
	app.GET("/genKey", getKey)
	app.GET("/logout", logout)

	app.GET("/login", renderLogin)
	app.GET("/signin", renderSingin)

	app.POST("/login", loginTry)
	app.POST("/signin", signinTry)

	app.GET("/apiTest", renderApiTest)
	app.GET("/apiTest/posteaza", renderApiTestPost)

	//Functiile din api.go
	app.GET("/api/haine", returneazaHaine)
	app.POST("/api/posteaza", posteaza)
	app.POST("/api/POST", posteazaHaine)
	app.GET("/api/haine:id", returneazaHainaDupaId)
	app.GET("/apiTest/haine:id", returneazaHainaDupaIdTest)

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
	sessionDelete(c)

	//redirectam catre pagina de login
	return c.Redirect(http.StatusSeeOther, "/login")
}

// Functie care randeaza pagina de testare a api ului
func renderApiTest(c echo.Context) error {
	//vom trimite toate hainele in pagina

	if !isLoggedIn(c) {
		c.Redirect(302, "/login")
		return nil
	}

	//incepem prin a face un rqeuest catre api pentru a lua toate hainele
	cheie := getCheieFromDB(c)
	if cheie == "" {
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return nil
	}
	linkPtReq := "http://localhost:8080/api/haine?key=" + cheie

	req, err := http.NewRequest("GET", linkPtReq, nil)
	check(err)

	cookie, err := c.Request().Cookie("session")
	check(err)

	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	body, err := io.ReadAll(resp.Body)
	check(err)
	resp.Body.Close()

	// Unmarshal the JSON response
	var haine []haina
	err = json.Unmarshal(body, &haine)
	check(err)

	// Acum putem sa folosim hainele
	return c.Render(http.StatusOK, "apiTest.html", map[string]interface{}{
		"haine": haine,
	})
}

// Functie care randeaza pagina de testare a postarii in api
func renderApiTestPost(c echo.Context) error {
	if !isLoggedIn(c) {
		c.Redirect(302, "/login")
		return nil
	}

	//randam pagina
	return c.Render(http.StatusOK, "apiTestPost.html", nil)

}

// Functie care da requestul de post catre api (cu cheie)
func posteaza(c echo.Context) error {
	if !isLoggedIn(c) {
		c.Redirect(302, "/login")
		return nil
	}

	cheie := getCheieFromDB(c)
	if cheie == "" {
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return nil
	}

	//preluam toate datele din formular
	nume := c.FormValue("Nume")
	tip := c.FormValue("Tip")
	culoare := c.FormValue("Culoare")
	marime := c.FormValue("Marime")
	sex := c.FormValue("Sex")
	pret := c.FormValue("Pret")

	//convertim tipul de date
	pretFloat, err := strconv.ParseFloat(pret, 32)
	check(err)
	tipInt, err := strconv.ParseInt(tip, 10, 64)
	check(err)

	var sexBool bool
	//verificam sexul
	if sex == "male" {
		sexBool = true
	} else if sex == "female" {
		sexBool = false
	}

	//creem haina
	haina := haina{
		Pret:    float32(pretFloat),
		Nume:    nume,
		Tip:     tipInt,
		Culoare: culoare,
		Marime:  marime,
		Sex:     sexBool,
	}

	//facem requestul catre api
	linkPtReq := "http://localhost:8080/api/POST?key=" + cheie

	js, err := json.Marshal(haina)
	check(err)

	req, err := http.NewRequest("POST", linkPtReq, bytes.NewBuffer(js))
	check(err)

	cookie, err := c.Request().Cookie("session")
	check(err)

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	//verificam daca a fost postat cu succes
	if resp.StatusCode == 201 {
		//daca da, trimitem un mesaj
		fmt.Println("Haina a fost postata cu succes")
	}

	//redirectam catre pagina de testare a api ului
	http.Redirect(c.Response(), c.Request(), "/apiTest", http.StatusSeeOther)
	return nil
}

// Functie care returneaza haina dupa id
func returneazaHainaDupaIdTest(c echo.Context) error {
	if !isLoggedIn(c) {
		c.Redirect(302, "/login")
		return nil
	}

	//preluam id ul din url
	id := c.Param("id")

	id = strings.Replace(id, ":", "", 1)

	//facem requestul catre api
	cheie := getCheieFromDB(c)
	if cheie == "" {
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return nil
	}

	linkPtReq := "http://localhost:8080/api/haine:" + id + "?key=" + cheie

	req, err := http.NewRequest("GET", linkPtReq, nil)
	check(err)

	cookie, err := c.Request().Cookie("session")
	check(err)

	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	body, err := io.ReadAll(resp.Body)
	check(err)
	resp.Body.Close()

	// Unmarshal the JSON response
	var haina haina
	err = json.Unmarshal(body, &haina)
	check(err)

	nume := haina.Nume
	culoare := haina.Culoare
	marime := haina.Marime
	pret := haina.Pret

	// Acum putem sa folosim haina
	return c.Render(http.StatusOK, "haina.html", map[string]interface{}{
		"Nume":    nume,
		"Culoare": culoare,
		"Marime":  marime,
		"Pret":    pret,
	})
}
