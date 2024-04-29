package apiFiles

import (
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

/*func MainHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("./HTML&CSS/index.html"))
	templ.Execute(w, nil)

	contor(w)
}

type idk struct {
	contor int64
}

func contor(w http.ResponseWriter) {
	var idkC idk
	idkC.contor = 0

	templ := template.Must(template.New("contor").Parse("{{.contor}}"))

	templ.Execute(w, idkC)
} */

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, contor string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, contor, data)
}
