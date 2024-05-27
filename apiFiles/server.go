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

type hainaCuId struct {
	Id      int     `json:"id"`
	Nume    string  `json:"nume"`
	Culoare string  `json:"culoare"`
	Marime  string  `json:"marime"`
	Pret    float32 `json:"pret"`
	Sex     bool    `json:"sex"`
	Tip     int64   `json:"tip"`
}

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
	app.GET("/apiTest/haina:id", returneazaHainaDupaIdTest)
	app.GET("/apiTest/filtrat", filter)
	app.POST("/api/filtreaza", filtreaza)

	app.DELETE("/apiTest/haine:id", stergeHaina)
	app.DELETE("/api/delete:id", stergeHainaDupaId)

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
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}

	cookie, err := c.Request().Cookie("session")
	if err != nil {
		fmt.Println("Eroare la preluarea cookie ului")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la request")
	}

	body, err := io.ReadAll(resp.Body)
	check(err)
	resp.Body.Close()

	// Unmarshal the JSON response
	var haine []haina
	err = json.Unmarshal(body, &haine)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare cu datele")
	}

	var iduri []int
	//preluam id urile din baza de date in ordine crescatoare
	db := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT id FROM haine")
	if err != nil {
		fmt.Println("Eroare la preluarea id urilor")
	}

	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		iduri = append(iduri, id)
	}

	var haineId []hainaCuId

	for i := 0; i < len(haine); i++ {
		haineId = append(haineId, hainaCuId{
			Id:      iduri[i],
			Nume:    haine[i].Nume,
			Pret:    haine[i].Pret,
			Sex:     haine[i].Sex,
			Tip:     haine[i].Tip,
			Culoare: haine[i].Culoare,
			Marime:  haine[i].Marime,
		})
	}

	// Acum putem sa folosim hainele
	return c.Render(http.StatusOK, "apiTest.html", map[string]interface{}{
		"haine": haineId,
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

	//veridicam daca sunt completate fieldurile
	if nume == "" || tip == "" || culoare == "" || marime == "" || sex == "" || pret == "" {
		return c.Redirect(http.StatusSeeOther, "/apiTest")
	}

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
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare cu datele")
	}

	req, err := http.NewRequest("POST", linkPtReq, bytes.NewBuffer(js))

	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la postare")
	}

	cookie, err := c.Request().Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Haina nu a putut fi postata")
	}

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
	if err != nil {
		fmt.Println("Eroare la request")
	}

	cookie, err := c.Request().Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la citirea raspunsului")
	}
	resp.Body.Close()

	// Unmarshal the JSON response
	var haina haina
	err = json.Unmarshal(body, &haina)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare cu datele")
	}

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

// Functie care sterge haina
func stergeHaina(c echo.Context) error {
	if !isLoggedIn(c) {
		c.Redirect(302, "/login")
		return nil
	}

	//preluam id ul din url(int ul de dupa :)
	id := c.Param("id")

	id = strings.Replace(id, ":", "", 1)

	//facem requestul catre api
	cheie := getCheieFromDB(c)
	if cheie == "" {
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return nil
	}

	linkPtReq := "http://localhost:8080/api/delete:" + id + "?key=" + cheie

	req, err := http.NewRequest("DELETE", linkPtReq, nil)

	if err != nil {
		fmt.Println("Haina nu poate fi stearsa")
	}

	cookie, err := c.Request().Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Haina nu a putut fi stearsa")
	}

	//verificam daca a fost stearsa cu succes
	if resp.StatusCode == 200 {
		//daca da, trimitem un mesaj
		fmt.Println("Haina a fost stearsa cu succes")
	}

	//redirectam catre pagina de testare a api ului
	http.Redirect(c.Response(), c.Request(), "/apiTest", http.StatusSeeOther)
	return nil
}

// Functie care filtreaza hainele
func filter(c echo.Context) error {
	if !isLoggedIn(c) {
		c.Redirect(302, "/login")
		return nil
	}

	//preluam datele din url
	marime := c.QueryParam("Marime")
	culoare := c.QueryParam("Culoare")
	tip := c.QueryParam("Tip")
	sex := c.QueryParam("Sex")
	pretMare := c.QueryParam("pretMare")
	pretMic := c.QueryParam("pretMic")

	//vom face un request catre api cu parametrii
	cheie := getCheieFromDB(c)
	if cheie == "" {
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return nil
	}

	linkPtReq := "http://localhost:8080/api/filtreaza?key=" + cheie

	//punem cookie un in request
	req, err := http.NewRequest("POST", linkPtReq, nil)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la request")
	}

	cookie, err := c.Request().Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	req.AddCookie(cookie)

	//veridicam ca toate sa nu fie goale

	if marime == "" && culoare == "" && tip == "" && sex == "" && pretMare == "" && pretMic == "" {
		http.Redirect(c.Response(), c.Request(), "/apiTest", http.StatusSeeOther)
		return nil
	}

	//punem parametrii in request
	q := req.URL.Query()
	if marime != "" {
		q.Add("marime", marime)
	}
	if culoare != "" {
		q.Add("culoare", culoare)
	}
	if sex != "" {
		q.Add("sex", sex)
	}

	if pretMare != "" {
		q.Add("pretMare", pretMare)
	}

	if pretMic != "" {
		q.Add("pretMic", pretMic)
	}

	if tip != "" {
		q.Add("tip", tip)
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare la citirea raspunsului")
	}

	resp.Body.Close()

	// Unmarshal the JSON response
	var haineFiltrate []haina
	err = json.Unmarshal(body, &haineFiltrate)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "Eroare cu datele")
	}

	//preluam id urile din baza de date DOAR ale hainelor filtrate, in oridinea filtrarii

	var haineId []hainaCuId

	for i := 0; i < len(haineFiltrate); i++ {
		db := connectToSQL()
		defer db.Close()

		//preluam numele si marimea
		nume := haineFiltrate[i].Nume
		marime := haineFiltrate[i].Marime
		culoare := haineFiltrate[i].Culoare
		tip := haineFiltrate[i].Tip

		var id int
		err = db.QueryRow("select id from haine where haina ->> '$.nume' = ? AND haina ->> '$.marime' = ? AND haina ->> '$.culoare' = ? AND haina ->> '$.tip' = ?", nume, marime, culoare, tip).Scan(&id)
		if err != nil {
			return c.String(400, "Eroare la preluarea id ului")
		}

		haineId = append(haineId, hainaCuId{
			Id:      id,
			Nume:    haineFiltrate[i].Nume,
			Pret:    haineFiltrate[i].Pret,
			Tip:     haineFiltrate[i].Tip,
			Sex:     haineFiltrate[i].Sex,
			Marime:  haineFiltrate[i].Marime,
			Culoare: haineFiltrate[i].Culoare,
		})
	}

	// Acum putem sa folosim hainele
	return c.Render(http.StatusOK, "apiTest.html", map[string]interface{}{
		"haine": haineId,
	})
}
