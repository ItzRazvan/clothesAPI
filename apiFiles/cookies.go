package apiFiles

/*func createCookie(c echo.Context, id int) error {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	//Stocam id ul userului in cookie
	cookie.Value = strconv.FormatInt(int64(id), 10)
	cookie.MaxAge = 60 * 60 * 24 * 3 // 1 week
	c.SetCookie(cookie)
	return nil
}

// Functie care verifica daca userul este logat
func isLoggedIn(c echo.Context) bool {
	_, err := c.Cookie("session")
	return err == nil
} */

// Functie care returneaza id ul userului din cookie (id ul este stocat acolo)
/*func getIdFromCookie(c echo.Context) int {
	if isLoggedIn(c) {
		cookie, err := c.Cookie("session")
		check(err)
		id, err := strconv.Atoi(cookie.Value)
		check(err)
		return id
	}
	return 0

} */

// Functie care sterge cookie ul cand userul se delogeaza
/*func removeCookie(c echo.Context) {
	cookie, err := c.Cookie("session")
	if err != nil {
		return
	}
	cookie.Name = "session"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-96 * time.Hour)
	cookie.MaxAge = -1
	c.SetCookie(cookie)
}
*/
