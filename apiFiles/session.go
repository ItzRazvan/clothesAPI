package apiFiles

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func getSessionKey() []byte {
	db := connectToSQL()
	defer db.Close()

	var key string
	err := db.QueryRow("SELECT cheie FROM sessionKey").Scan(&key)
	if err != nil {
		key = genereazaStringRandom(32)
		_, err = db.Exec("INSERT INTO sessionKey (cheie) VALUES (?)", key)
		if err != nil {
			fmt.Println("Eroare la inserarea cheii in baza de date")
			return nil
		}
	}

	return []byte(key)
}

var (
	key   = getSessionKey()
	store = sessions.NewCookieStore(key)
)

func isLoggedIn(c echo.Context) bool {
	session, _ := store.Get(c.Request(), "session")

	auth, ok := session.Values["authenticated"].(bool)

	return auth && ok

}

func sessionInit(c echo.Context, id int) error {
	session, err := store.Get(c.Request(), "session")

	if err != nil {
		fmt.Println("Eroare la preluarea sesiunii")
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return err
	}

	//setam userul ca authentificat
	session.Values["authenticated"] = true

	//setam id ul userului
	session.Values["id"] = id

	//Facem setarile sesiunii
	session.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24 * 3,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	err = session.Save(c.Request(), c.Response())

	if err != nil {
		fmt.Println("Eroare la salvarea sesiunii")
		http.Redirect(c.Response(), c.Request(), "/login", http.StatusSeeOther)
		return err
	}

	return nil
}

func sessionDelete(c echo.Context) {
	session, err := store.Get(c.Request(), "session")
	if err != nil {
		return
	}
	session.Values["authenticated"] = false
	session.Save(c.Request(), c.Response())

}

func getIdFromSession(c echo.Context) int {
	session, err := store.Get(c.Request(), "session")
	if err != nil {
		fmt.Println("Eroare la preluarea sesiunii")
		c.Redirect(http.StatusSeeOther, "/login")
		return 0
	}

	var id int
	if session.IsNew {
		fmt.Println("Sesiunea nu exista")

		//Daca nu exista sesiunea, vom lua id ul din cookieul cu numele session

		cookie, err := c.Cookie("session")
		if err != nil {
			fmt.Println("Cookieul nu exista")
			c.Redirect(http.StatusSeeOther, "/login")
			return 0
		}

		id, err = strconv.Atoi(cookie.Value)
		if err != nil {
			fmt.Println("Eroare la convertirea id ului din cookie")
			c.Redirect(http.StatusSeeOther, "/login")
			return 0
		}

		return id
	}

	id = session.Values["id"].(int)
	return id
}
