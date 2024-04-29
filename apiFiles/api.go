package apiFiles

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

//Un API folosit pentru a adauga haine intr-un shop online sau fizic

type haina struct {
	ID      int64   `json:"id"`
	Pret    float64 `json:"pret"`
	Nume    string  `json:"nume"`
	Tip     int64   `json:"tip"` // 1 - hanorac, 2 - blugi, 3 - pantaloni cargo, 4 - tricou
	Culoare string  `json:"culoare"`
	Marime  string  `json:"marime"` // XS,S,M,L,XL,XXL
	Sex     bool    `json:"sex"`    // 0 - Femei, 1 - Barbati
}

//Date de inceput

var haine = []haina{
	{ID: 1, Pret: 69.99, Nume: "ZW COLLECTION BOOTCUT MID-RISE CONTOUR JEANS", Tip: 2, Culoare: "Albastru", Marime: "L", Sex: false},
	{ID: 2, Pret: 129.99, Nume: "JEWEL NECKLACE HOODIE", Tip: 1, Culoare: "Rosu", Marime: "S", Sex: true},
}

// Functie care verifica sa nu fie vreo eroare
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Functie care returneaza hainele ca un JSON
func returneazaHaine(c echo.Context) {
	c.JSON(http.StatusOK, haine)
}

// functie pentru a posta o haina
func posteazaHaine(c echo.Context) error {
	var hainaNoua haina

	err := c.Bind(&hainaNoua)
	check(err)

	haine = append(haine, hainaNoua)

	return c.JSON(http.StatusCreated, hainaNoua)
}

// functie care returneaza o haina dupa id-ul lui
func returneazaHainaDupaId(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	check(err)
	for _, haina := range haine {
		if haina.ID == idInt {
			c.JSON(http.StatusOK, haina)
			return nil
		}
	}
	return echo.ErrBadRequest
}

/* func main() {
	router := gin.Default()

	router.GET("/haine", returneazaHaine)
	router.GET("/haina/:id", returneazaHainaDupaId)

	router.POST("/posteaza", posteazaHaine)

	err := router.Run("localhost:8080")
	check(err)

} */
