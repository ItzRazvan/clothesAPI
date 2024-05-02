package apiFiles

import (
	"fmt"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/labstack/echo/v4"
	passwordValidator "github.com/wagslane/go-password-validator"
)

var verifier = emailverifier.NewVerifier()

// Functie care verifica daca datele de signin sunt corecte, emailul nu exista deja in baza de date
func signinTry(c echo.Context) error {
	//preluam datele din formular
	email := c.FormValue("email")
	parola := c.FormValue("parola")

	fmt.Println(email, parola)

	//verificam daca parola este destul de puternica
	const minEntropy float64 = 70
	err := passwordValidator.Validate(parola, minEntropy)
	if err != nil {
		return c.String(400, "Parola prea slaba")
	}

	parolaHash := generateToken(parola)

	//verificam daca emailul este valid
	ret, err := verifier.Verify(email)
	if err != nil {
		return c.String(400, "Email invalid")
	}
	if !ret.Syntax.Valid {
		return c.String(400, "Sintaxa invalida")
	}

	//verificam daca emailul exista deja in baza de date
	db := connectToSQL()
	defer db.Close()

	var emailDinDB string
	err = db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&emailDinDB)

	if err == nil {
		return c.String(400, "Emailul exista deja in baza de date")
	}

	//stergem indexul existent de pe email
	_, _ = db.Exec("DROP INDEX email ON users")

	//daca emailul nu exista, il adaugam in baza de date
	_, err = db.Exec("INSERT INTO users (email, parola) VALUES (?, ?)", email, parolaHash)
	check(err)

	//cream un index nou pentru email
	_, err = db.Exec("CREATE UNIQUE INDEX email ON users (email)")
	check(err)

	//daca totul este ok, cream un cookie pentru a tine minte ca userul este logat
	createCookie(c)

	//daca totul este ok, returnam un mesaj de succes
	return c.String(200, "Inregistrare cu succes")
}
