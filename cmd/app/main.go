package main

import (
	"ash_cheatsheet/internal/db"
	"ash_cheatsheet/internal/entities"
	"github.com/pav5000/go-common/lambda"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type SectionData struct {
	SectionName string
	Cards       []CardView
}

type AddCardData struct {
	SectionName string
}

type CardView struct {
	Name        string
	Description string
}

func main() {
	tmpl, err := template.ParseGlob("templates/*.htm")
	if err != nil {
		panic(err)
	}

	dbsvc, err := db.New()
	if err != nil {
		panic(err)
	}

	sectionHtm := tmpl.Lookup("section.htm")
	if sectionHtm == nil {
		panic("section htm is nil")
	}

	cardHtm := tmpl.Lookup("add_card.htm")
	if cardHtm == nil {
		panic("new_card htm is nil")
	}

	http.HandleFunc("GET /static/{filename}",
		func(writer http.ResponseWriter, request *http.Request) {
			fileName := request.PathValue("filename")

			switch {
			case strings.HasSuffix(fileName, ".css"):
				writer.Header().Add("Content-Type", "text/css")
			case strings.HasSuffix(fileName, ".woff2"):
				writer.Header().Add("Content-Type", "application/font-woff2")
			default:
				panic("invalid file")
			}

			f, err := os.Open("static/" + fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := io.Copy(writer, f); err != nil {
				panic(err)
			}
		},
	)

	http.HandleFunc("GET /sheet/{section}/add",
		func(writer http.ResponseWriter, request *http.Request) {
			section := request.PathValue("section")

			if err := cardHtm.Execute(writer, AddCardData{SectionName: section}); err != nil {
				panic(err)
			}
		},
	)

	http.HandleFunc("GET /sheet/{section}",
		func(writer http.ResponseWriter, request *http.Request) {
			section := request.PathValue("section")

			cards, err := dbsvc.SelectAllCardsBySection(request.Context(), section)
			if err != nil {
				panic(err)
			}

			if err := sectionHtm.Execute(writer, SectionData{
				SectionName: section,
				Cards: lambda.MapSlice(cards, func(card entities.Card) CardView {
					return CardView{
						Name:        card.Name,
						Description: card.Description,
					}
				}),
			}); err != nil {
				panic(err)
			}
		},
	)

	err = http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		panic(err)
	}
}
