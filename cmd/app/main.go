package main

import (
	"ash_cheatsheet/internal/db"
	"ash_cheatsheet/internal/entities"
	"github.com/pav5000/go-common/lambda"
	"html/template"
	"net/http"
)

type SectionData struct {
	SectionName string
	Cards       []CardView
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