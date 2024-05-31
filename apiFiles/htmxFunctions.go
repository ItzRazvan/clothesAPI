package apiFiles

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Structura care contine template-urile
type Template struct {
	templates *template.Template
}

// Functie care da rander la html cu template
func (t *Template) Render(w io.Writer, index string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, index, data)
}

// Genereaza un string care urmeaza sa fie hashuit
func genereazaStringRandom(size int) string {
	//caracterele folosite
	const caractere = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+{}:<>?,./;[]'-=`~"

	a := make([]byte, size)

	//generam un string random
	for i := range a {
		a[i] = caractere[rand.Intn(len(caractere))]
	}

	return string(a)
}

// Genereaza o cheie random
func genereazaCheieRandom() string {
	//generam un string random
	randString := genereazaStringRandom(25)
	//hashuim stringul
	hashedString := generateToken(randString)

	return hashedString
}

// Trecem un string random printr o functie de hash pt a avea o cheie cat mai unica
func generateToken(text string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return text
	}

	return string(hash)
}

// Functie care schimba cheia
func schimbareCheie(c echo.Context, id int) {
	db := connectToSQL()
	defer db.Close()

	cheie := genereazaCheieRandom()

	//daca id ul este 0 atunci luam id ul din sesiune
	if id == 0 {
		idBun := getIdFromSession(c)
		id = idBun

		//daca id ul este 0 inseamna ca nu este nimeni logat
		if id == 0 {
			c.Redirect(302, "/login")
			return
		}
	}

	//schimbam cheia in baza de date
	_, err := db.Exec("UPDATE users SET cheie = ? WHERE id = ?", cheie, id)
	if err != nil {
		fmt.Fprintf(c.Response(), "Eroare la schimbarea cheii")
	}
}

// fucntie care ia cheia din baza de date
func getCheieFromDB(c echo.Context) string {
	db := connectToSQL()
	defer db.Close()

	//luam id ul din sesiune
	id := getIdFromSession(c)
	if id == 0 {
		fmt.Println("Userul nu este detectat logat")
		return ""
	}

	//selectam cheia din baza de date
	var cheie string
	err := db.QueryRow("SELECT cheie FROM users WHERE id = ?", id).Scan(&cheie)
	if err != nil {
		fmt.Fprintf(c.Response(), "Eroare la preluarea cheii")
		return ""
	}

	return cheie
}

// Functie care te conecteaza la baza de date
func connectToSQL() *sql.DB {
	db, err := sql.Open("mysql", "root:razvan@tcp(127.0.0.1:3306)/clothesAPI")
	if err != nil {
		log.Fatal("Eroare la conectarea la baza de date")
	}
	db.Ping()
	return db
}
