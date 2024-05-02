package apiFiles

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Functie care verifica daca datele de login sunt corecte
func loginTry(c echo.Context) error {
	//preluam datele din formular
	email := c.FormValue("email")
	parola := c.FormValue("parola")

	//verificam daca datele sunt corecte
	db := connectToSQL()
	defer db.Close()

	var emailDinDB string
	var parolaDinDB string
	err := db.QueryRow("SELECT email, parola FROM users WHERE email = ?", email).Scan(&emailDinDB, &parolaDinDB)

	if err != nil {
		return c.String(400, "Emailul sau parola sunt gresite")
	}

	//verificam daca parola este corecta
	err = bcrypt.CompareHashAndPassword([]byte(parolaDinDB), []byte(parola))
	if err != nil {
		return c.String(400, "Emailul sau parola sunt gresite")
	}

	//daca totul e corect, crean un cookie pentru a tine minte ca userul este logat
	createCookie(c)

	//Returnam un mesaj de succes
	return c.String(200, "Logare cu succes")
}