package main

import (
	"ash_cheatsheet/internal/cards"
	"ash_cheatsheet/internal/db"
	"ash_cheatsheet/internal/entities"
	"ash_cheatsheet/internal/errs"
	"ash_cheatsheet/internal/handlers/preview"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pav5000/go-common/lambda"
	"github.com/pav5000/logger"
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
	Id          int64
}

type AddCardResult struct {
	Message string
	Class   string // is-danger
}

type EditCardResult struct {
	Message string
	Class   string
}

type EditCardPage struct {
	Id          int64
	Name        string
	Description string
	SectionName string
}

func main() {
	tmpl, err := template.ParseGlob("templates/*.htm")
	if err != nil {
		logger.Fatal(err.Error())
	}

	dbsvc, err := db.New()
	if err != nil {
		logger.Fatal(err.Error())
	}

	sectionHtm := tmpl.Lookup("section.htm")
	if sectionHtm == nil {
		panic("section htm is nil")
	}

	cardHtm := tmpl.Lookup("add_card.htm")
	if cardHtm == nil {
		panic("new_card htm is nil")
	}

	editCardHtm := tmpl.Lookup("edit.htm")
	if editCardHtm == nil {
		panic("edit card htm is nil")
	}

	editCardResultHtm := tmpl.Lookup("edit_result.htm")
	if editCardResultHtm == nil {
		panic("edit_result htm is nil")
	}

	addResult := tmpl.Lookup("add_result.htm")
	if addResult == nil {
		panic("add result htm is nil")
	}

	cardsSrv := cards.New(dbsvc)

	http.HandleFunc("GET /static/{filename}",
		func(writer http.ResponseWriter, request *http.Request) {
			fileName := request.PathValue("filename")

			switch {
			case strings.HasSuffix(fileName, ".css"):
				writer.Header().Add("Content-Type", "text/css")
			case strings.HasSuffix(fileName, ".woff2"):
				writer.Header().Add("Content-Type", "application/font-woff2")
			case strings.HasSuffix(fileName, ".js"):
				writer.Header().Add("Content-Type", "text/javascript")
			default:
				panic("invalid file")
			}

			f, err := os.Open("static/" + fileName)
			if err != nil {
				logger.Fatal(err.Error())
			}
			defer f.Close()

			if _, err := io.Copy(writer, f); err != nil {
				logger.Fatal(err.Error())
			}
		},
	)

	http.HandleFunc("POST /preview", preview.New())

	http.HandleFunc("GET /sheet/{section}/add",
		func(writer http.ResponseWriter, request *http.Request) {
			section := request.PathValue("section")

			if err := cardHtm.Execute(writer, AddCardData{SectionName: section}); err != nil {
				logger.Fatal(err.Error())
			}
		},
	)

	http.HandleFunc("GET /sheet/{section}/{card_id}/edit",
		func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("card_id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				logger.Fatal(err.Error())
			}

			card, err := cardsSrv.GetCardByID(int64(id))
			if err != nil {
				logger.Fatal(err.Error())
			}

			if err := editCardHtm.Execute(w, EditCardPage{
				Id:          card.Id,
				Name:        card.Name,
				Description: card.Description,
				SectionName: card.Section,
			}); err != nil {
				logger.Fatal(err.Error())
			}
		},
	)

	http.HandleFunc("POST /sheet/{section}/{card_id}/edit", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("card_id")
		cardId, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Fatal(err.Error())
		}

		rawFormData, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Fatal(err.Error())
		}

		vals, err := url.ParseQuery(string(rawFormData))
		if err != nil {
			logger.Fatal(err.Error())
		}

		name := vals.Get("name")
		description := vals.Get("description")

		if err := cardsSrv.UpdateCardByID(int64(cardId), name, description); err != nil {
			logger.Fatal(err.Error())
		}

		if err := editCardResultHtm.Execute(w, EditCardResult{Message: "Карточка изменена"}); err != nil {
			logger.Fatal(err.Error())
		}
	})

	http.HandleFunc("DELETE /sheet/{section}/{card_id}",
		func(writer http.ResponseWriter, request *http.Request) {
			sectionName := request.PathValue("section")
			idStr := request.PathValue("card_id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				logger.Fatal(err.Error())
			}

			if err := cardsSrv.DeleteCard(int64(id), sectionName); err != nil {
				logger.Fatal(err.Error())
			}
		},
	)

	http.HandleFunc("POST /sheet/{section}/add",
		func(writer http.ResponseWriter, request *http.Request) {
			rawFormData, err := io.ReadAll(request.Body)
			if err != nil {
				logger.Fatal(err.Error())
			}

			vals, err := url.ParseQuery(string(rawFormData))
			if err != nil {
				logger.Fatal(err.Error())
			}

			sectionName := request.PathValue("section")

			var res AddCardResult

			err = cardsSrv.CreateNewCard(vals.Get("name"), vals.Get("description"), sectionName)
			switch {
			case errors.Is(err, errs.ErrEmptyCardName):
				res = AddCardResult{Message: "Имя карточки не может быть пустым!", Class: "is-danger"}
			case errors.Is(err, errs.ErrEmptyCardDesc):
				res = AddCardResult{Message: "Описание карточки не может быть пустым!", Class: "is-danger"}
			default:
				res = AddCardResult{Message: "Новая карточка добавлена в секцию " + sectionName}
			}
			if err := addResult.Execute(writer, res); err != nil {
				logger.Fatal(err.Error())
			}
		},
	)

	http.HandleFunc("GET /sheet/{section}",
		func(writer http.ResponseWriter, request *http.Request) {
			section := request.PathValue("section")

			cardsBySection, err := dbsvc.SelectAllCardsBySection(request.Context(), section)
			if err != nil {
				logger.Fatal(err.Error())
			}

			if err := sectionHtm.Execute(writer, SectionData{
				SectionName: section,
				Cards: lambda.MapSlice(cardsBySection, func(card entities.Card) CardView {
					return CardView{
						Name:        card.Name,
						Description: card.Description,
						Id:          card.Id,
					}
				}),
			}); err != nil {
				logger.Fatal(err.Error())
			}
		},
	)

	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
