package apiFiles

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func createCookie(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Expires = time.Now().Add(96 * time.Hour)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "cookie set")
}

func isLoggedIn(c echo.Context) bool {
	_, err := c.Cookie("session")
	return err == nil
}
