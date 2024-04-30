package apiFiles

import (
	"database/sql"
	"io"
	"log"
	"math/rand"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, index string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, index, data)
}

// Genereaza un string care urmeaza sa fie hashuit
func genereazaStringRandom(size int) string {
	//caracterele folosite
	const caractere = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+{}:<>?,./;[]'-=`~"

	a := make([]byte, size)

	for i := range a {
		a[i] = caractere[rand.Intn(len(caractere))]
	}

	return string(a)
}

func genereazaCheieRandom() string {
	//generate an random string of size 25
	randString := genereazaStringRandom(25)
	hashedString := generateToken(randString)

	return hashedString
}

// Trecem un string random printr o functie de hash pt a avea o cheie cat mai unica
func generateToken(text string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

// Functia verifica daca cheia este generata, pentru a nu genera inca una, doar in caz ca unserul apasa pe butonul de generare
func existaSauGenereazaCheie() string {
	db, err := sql.Open("mysql", "root:razvan2007@tcp(127.0.0.1:3306)/clothesAPI")
	check(err)
	err = db.Ping()
	check(err)

	//verificam daca exista deja o cheie in baza de date
	var cheieDinDB string
	err = db.QueryRow("SELECT value FROM cheie").Scan(&cheieDinDB)

	//daca exista, returnam cheia
	if err == nil {
		return cheieDinDB
	} else {
		//daca nu exista, generam una noua
		cheieNoua := genereazaCheieRandom()
		_, err = db.Exec("INSERT INTO cheie (value) VALUES (?)", cheieNoua)
		check(err)
		return cheieNoua
	}
}

func schimbareCheie() {
	db, err := sql.Open("mysql", "root:razvan2007@tcp(127.0.0.1:3306)/clothesAPI")
	check(err)
	err = db.Ping()
	check(err)

	cheie := genereazaCheieRandom()

	_, err = db.Exec("UPDATE cheie SET value = ?", cheie)
	check(err)
}
