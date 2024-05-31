package apiFiles

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Un API folosit pentru a adauga haine intr-un shop online sau fizic
// Pentru a face requesturi la api, vei folosi apiKey ul generat pe site ul nostru:
// (vei face req la adresa http://site/api/haine?key=...)
// Pentru a sterge o haina req la adresa http://site/api/delete:id?key=...
// Pentru a posta o haina, vei face un request POST la adresa http://site/api/POST?key=...
// Pentru a vedea o haina dupa id, vei face un request GET la adresa http://site/api/haina/:id?key=...
// Vei avea nevoie de un struct haina, deoarece api ul va returna hainele in forma de struct din baza de date
type haina struct {
	Pret    float32 `json:"pret"`
	Nume    string  `json:"nume"`
	Tip     int64   `json:"tip"` // 1 - hanorac, 2 - blugi, 3 - pantaloni cargo, 4 - tricou
	Culoare string  `json:"culoare"`
	Marime  string  `json:"marime"` // XS,S,M,L,XL,XXL
	Sex     bool    `json:"sex"`    // 0 - Femei, 1 - Barbati
}

// functie cara adauga o haina in baza de date
func adaugaInBazaDeDate(hainaDeAdaugat haina) {
	//adauga hainele din var haine in baza de date
	db := connectToSQL()
	defer db.Close()

	//formatam haina in JSON
	j, err := json.Marshal(hainaDeAdaugat)
	if err != nil {
		fmt.Println("Eraore la adaugar hainei in baza de date")
		return
	}

	//Scoatem indexu daca exista, ca sa il refacem apoi
	//Nu bagam in seama eraorea pt ca nu conteaza daca nu are index, il vom crea oricum
	_, _ = db.Exec("DROP INDEX indexHaine ON haine")

	//inseram haina in baza de date
	_, err = db.Exec("INSERT INTO haine (haina) VALUES (?)", j)
	if err != nil {
		fmt.Println("Eroare la adaugare hainei in baza de date")
		return
	}

	//creem indexul
	_, err = db.Exec("CREATE UNIQUE INDEX indexHaine ON haine (id, haina)")
	if err != nil {
		fmt.Println("Eroare la crearea indexului")
	}
}

// Functie care verifica sa nu fie vreo eroare si opreste programul daca este cazul
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Functie care returneaza hainele
func returneazaHaine(c echo.Context) error {
	//daca cheia este buna, returneaza hainele
	if keyIsOk(c) {
		haine := iaHaineDinBazaDeDate()
		//returneaza hainele in format json
		return c.JSON(http.StatusOK, haine)
	}
	return echo.ErrBadRequest
}

// functie pentru a posta o haina
func posteazaHaine(c echo.Context) error {
	if keyIsOk(c) {
		var hainaNoua haina

		//luam haina din request
		err := c.Bind(&hainaNoua)
		if err != nil {
			return echo.ErrBadRequest
		}

		//le adaugam in baza de date
		adaugaInBazaDeDate(hainaNoua)

		return c.JSON(http.StatusCreated, hainaNoua)
	}
	return echo.ErrBadRequest
}

// functie care sterge o haina dupa id
func stergeHainaDupaId(c echo.Context) error {
	if keyIsOk(c) {
		//luam id ul din parametrii url
		id := c.Param("id")
		//stergem : din id
		id = strings.Replace(id, ":", "", 1)
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}
		db := connectToSQL()
		defer db.Close()

		//Selectam al idInt-n id
		var idDB int
		err = db.QueryRow("SELECT id FROM haine order by id limit ?,1", idInt-1).Scan(&idDB)
		if err != nil {
			return echo.ErrBadRequest
		}

		//il stergem din baza de date
		_, err = db.Exec("DELETE FROM haine WHERE id = ?", idDB)

		if err == nil {
			fmt.Println("Eroare la stergerea hainei din baza de date")
			return echo.ErrBadRequest
		}

		return c.String(http.StatusOK, fmt.Sprintf("Haina cu id ul %d a fost stearsa", idInt))
	}
	return echo.ErrBadRequest
}

// functie care returneaza o haina dupa id-ul lui
func returneazaHainaDupaId(c echo.Context) error {
	if keyIsOk(c) {
		//luam id ul din parametrii url
		id := c.Param("id")
		//stergem : din id
		id = strings.Replace(id, ":", "", 1)
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}

		//selectam din baza de date haine cu id ul idInt
		db := connectToSQL()
		defer db.Close()

		//selectam haina
		var hainaJson string
		err = db.QueryRow("SELECT haina FROM haine where id = ?", idInt).Scan(&hainaJson)

		if err != nil {
			return echo.ErrBadRequest
		}

		//formatam haina in structura haina
		var h haina
		err = json.Unmarshal([]byte(hainaJson), &h)

		if err != nil {
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, h)

	}
	return echo.ErrBadRequest
}

// functie care ia hainele din baza de date
func iaHaineDinBazaDeDate() []haina {
	var haine []haina

	db := connectToSQL()
	defer db.Close()

	//selectam hainele
	rows, err := db.Query("SELECT haina FROM haine")
	if err != nil {
		fmt.Println("Eroare la preluarea hainelor din baza de date")
		return nil
	}

	//parcurgem hainele si le adaugam in slice ul haine
	for rows.Next() {
		var h haina
		var hainaJson string
		err = rows.Scan(&hainaJson)
		if err != nil {
			continue
		}
		err = json.Unmarshal([]byte(hainaJson), &h)
		if err != nil {
			continue
		}
		haine = append(haine, h)
	}
	return haine
}

// Functie care filtreaza hainele
func filtreaza(c echo.Context) error {
	if keyIsOk(c) {
		//luam criteriile de filtrare din url
		tip := c.QueryParam("tip")
		culoare := c.QueryParam("culoare")
		marime := c.QueryParam("marime")
		sex := c.QueryParam("sex")

		sexBool := false
		if sex == "male" {
			sexBool = true
		}

		pretMare := c.QueryParam("pretMare")
		pretMic := c.QueryParam("pretMic")

		haine := iaHaineDinBazaDeDate()

		var haineFiltrate []haina

		//parcurgem hainele si le adaugam in slice ul haineFiltrate daca respecta criteriile
		for _, h := range haine {
			//excludem criterile care sunt ""

			if tip != "" {
				tipInt, err := strconv.ParseInt(tip, 10, 64)
				if err != nil {
					continue
				}
				if h.Tip != tipInt {
					continue
				}
			}

			if culoare != "" {
				if h.Culoare != culoare {
					continue
				}
			}

			if marime != "" {
				if h.Marime != marime {
					continue
				}
			}

			if sex != "" {
				if h.Sex != sexBool {
					continue
				}
			}

			if pretMare != "" {
				pretMareFloat, err := strconv.ParseFloat(pretMare, 32)
				if err != nil {
					continue
				}
				if float32(pretMareFloat) > h.Pret {
					continue
				}
			}

			if pretMic != "" {
				pretMicFloat, err := strconv.ParseFloat(pretMic, 32)
				if err != nil {
					continue
				}
				if float32(pretMicFloat) < h.Pret {
					continue
				}
			}

			haineFiltrate = append(haineFiltrate, h)
		}
		return c.JSON(http.StatusOK, haineFiltrate)
	}

	return nil
}

// Verificam daca cheia din url este buna
func keyIsOk(c echo.Context) bool {
	//luam cheia din url
	cheie := c.QueryParam("key")
	//luam cheia din baza de date
	cheieDinDB := getCheieFromDB(c)
	//verificam daca sunt egale
	return cheie == cheieDinDB
}
