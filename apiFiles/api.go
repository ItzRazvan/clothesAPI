package apiFiles

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Un API folosit pentru a adauga haine intr-un shop online sau fizic
// Pentru a face requesturi la api, vei folosi apiKey ul generat pe site ul nostru:
// (vei face req la adresa http://site/api/haine?key=...)
// Pentru a posta o haina, vei face un request POST la adresa http://site/api/POST?key=...
// Pentru a vedea o haina dupa id, vei face un request GET la adresa http://site/api/haina/:id?key=...
// Vei avea nevoie de un struct haina, deoarece api ul va returna hainele in forma de struct din baza de date
type haina struct {
	Pret    float64 `json:"pret"`
	Nume    string  `json:"nume"`
	Tip     int64   `json:"tip"` // 1 - hanorac, 2 - blugi, 3 - pantaloni cargo, 4 - tricou
	Culoare string  `json:"culoare"`
	Marime  string  `json:"marime"` // XS,S,M,L,XL,XXL
	Sex     bool    `json:"sex"`    // 0 - Femei, 1 - Barbati
}

//Date de inceput
//var haina1 = haina{Pret: 69.99, Nume: "ZW COLLECTION BOOTCUT MID-RISE CONTOUR JEANS", Tip: 2, Culoare: "Albastru", Marime: "L", Sex: false}
//var haina2 = haina{Pret: 129.99, Nume: "JEWEL NECKLACE HOODIE", Tip: 1, Culoare: "Rosu", Marime: "S", Sex: true}

func adaugaInBazaDeDate(hainaDeAdaugat haina) {
	//adauga hainele din var haine in baza de date
	db := connectToSQL()
	defer db.Close()

	j, err := json.Marshal(hainaDeAdaugat)
	check(err)

	_, err = db.Exec("INSERT INTO haine (haina) VALUES (?)", j)
	check(err)
}

// Functie care verifica sa nu fie vreo eroare
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Functie care returneaza hainele ca un un slice de struct haina
func returneazaHaine(c echo.Context) error {
	//daca cheia este buna, returneaza hainele
	if keyIsOk(c) {
		haine := iaHaineDinBazaDeDate()
		//returneaza hainele in format json
		return c.JSON(http.StatusOK, haine)
	}
	return echo.ErrBadRequest
}

// functie pentru a posta o haina
func posteazaHaine(c echo.Context) error {
	if !keyIsOk(c) {
		var hainaNoua haina

		err := c.Bind(&hainaNoua)
		check(err)

		adaugaInBazaDeDate(hainaNoua)

		return c.JSON(http.StatusCreated, hainaNoua)
	}
	return echo.ErrBadRequest
}

// functie care returneaza o haina dupa id-ul lui
func returneazaHainaDupaId(c echo.Context) error {
	if keyIsOk(c) {
		id := c.Param("id")
		//delete the : from id
		id = strings.Replace(id, ":", "", 1)
		idInt, err := strconv.ParseInt(id, 10, 64)
		check(err)
		db := connectToSQL()
		defer db.Close()

		var haina haina
		//luam haina careia id este idInt
		rows := db.QueryRow("SELECT haina FROM haine WHERE id = ?", idInt)
		var hainaJson string
		err = rows.Scan(&hainaJson)
		if err == nil {
			err = json.Unmarshal([]byte(hainaJson), &haina)
			if err == nil {
				return c.JSON(http.StatusOK, haina)
			}
			return c.String(http.StatusNotFound, fmt.Sprintf("Haina cu id ul %d nu a fost gasita", idInt))
		}
		return c.String(http.StatusNotFound, fmt.Sprintf("Haina cu id ul %d nu a fost gasita", idInt))

	}
	return echo.ErrBadRequest
}

// functie care ia hainele din baza de date
func iaHaineDinBazaDeDate() []haina {
	var haine []haina

	db := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT haina FROM haine")
	check(err)

	for rows.Next() {
		var h haina
		var hainaJson string
		err = rows.Scan(&hainaJson)
		check(err)
		err = json.Unmarshal([]byte(hainaJson), &h)
		check(err)
		haine = append(haine, h)
	}
	return haine
}

// Functie care porneste serverul
//In acest demo nu este folosita pt ca deja exista un server folosit
//Dar, voi folosi fix aceste functii in serverul folosit deja
/*func StartApi() {
	app := echo.New()

	app.GET("/api/haine", returneazaHaine)
	app.POST("/api/POST", posteazaHaine)
	app.GET("/api/haine/:id", returneazaHainaDupaId)

	app.Logger.Fatal(app.Start(":8080"))
}

*/

// Verificam daca cheia din url este buna
func keyIsOk(c echo.Context) bool {
	cheie := c.QueryParam("key")
	cheieDinDB := getCheieFromDB(c)
	return cheie == cheieDinDB
}
