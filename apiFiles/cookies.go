package apiFiles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func createCookie(c echo.Context, id int) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	//Stocam id ul userului in cookie
	cookie.Value = strconv.FormatInt(int64(id), 10)
	cookie.Expires = time.Now().Add(96 * time.Hour)
	c.SetCookie(cookie)
	return nil
}

// Functie care verifica daca userul este logat
func isLoggedIn(c echo.Context) bool {
	_, err := c.Cookie("session")
	return err == nil
}

// Functie care returneaza id ul userului din cookie (id ul este stocat acolo)
func getIdFromCookie(c echo.Context) int {
	cookie, err := c.Cookie("session")
	check(err)
	id, _ := strconv.Atoi(cookie.Value)
	return id
}

// Functie care sterge cookie ul cand userul se delogeaza
func removeCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-96 * time.Hour)
	cookie.MaxAge = -1
	c.SetCookie(cookie)
}
