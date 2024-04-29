package apiFiles

import (
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

// porneste serverului cu toate functiile sale
func ServerStart() {
	app := echo.New()

	//initializare template randerer
	t := &Template{
		templates: template.Must(template.ParseGlob("./views/index.html")),
	}
	app.Renderer = t
	app.GET("/", renderIndex)
	app.GET("/genKey", generateKey)

	app.Logger.Fatal(app.Start(":8080"))
}

// Afiseaza index.html cand este accesata pagina principala
func renderIndex(c echo.Context) error {

	cheie := existaSauGenereazaCheie()

	key := map[string]interface{}{
		"Cheie": cheie,
	}

	return c.Render(http.StatusOK, "apiKey", key)
}

// Genereaza o cheie noua dupa ce butonul de generare este apasat
func generateKey(c echo.Context) error {
	schimbareCheie()

	return nil
}

// onclick="setTimeout(()=>(window.location = window.location.href), 50)"
