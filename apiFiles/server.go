package apiFiles

import (
	"net/http"

	"github.com/fatih/color"
)

// initializarea serverului cu toate functiile sale
func ServerInit() {
	//initializare pagina principala
	serverFile := http.FileServer(http.Dir("../HTML&CSS"))
	http.Handle("/", serverFile)

	color.Green("\nServerul a pornit...\n")

	err := http.ListenAndServe(":8080", nil)
	check(err)
}
