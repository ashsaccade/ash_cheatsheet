package main

import (
	"ash_cheatsheet/internal/cards"
	"ash_cheatsheet/internal/db"
	"ash_cheatsheet/internal/handlers/addcard"
	"ash_cheatsheet/internal/handlers/deletecard"
	"ash_cheatsheet/internal/handlers/getadd"
	"ash_cheatsheet/internal/handlers/geteditcard"
	"ash_cheatsheet/internal/handlers/getsection"
	"ash_cheatsheet/internal/handlers/posteditcard"
	"ash_cheatsheet/internal/handlers/preview"
	"ash_cheatsheet/internal/handlers/static"
	"github.com/pav5000/logger"
	"html/template"
	"net/http"
)

func main() {
	tmpl, err := template.ParseGlob("templates/*.htm")
	if err != nil {
		logger.Fatal(err.Error())
	}

	dbsvc, err := db.New()
	if err != nil {
		logger.Fatal(err.Error())
	}

	cardsSrv := cards.New(dbsvc)

	editCardHandler := geteditcard.New(tmpl, cardsSrv)
	getAddCardHandler := getadd.New(tmpl)
	addCardHandler := addcard.New(cardsSrv, tmpl)
	postEditCard := posteditcard.New(tmpl, cardsSrv)
	deleteCard := deletecard.New(cardsSrv)
	sectionViewHandler := getsection.New(tmpl, cardsSrv)

	http.HandleFunc("GET /static/{filename}", static.New())
	http.HandleFunc("POST /preview", preview.New())
	http.HandleFunc("GET /sheet/{section}/add", getAddCardHandler.Handle())
	http.HandleFunc("GET /sheet/{section}/{card_id}/edit", editCardHandler.Handle())
	http.HandleFunc("POST /sheet/{section}/{card_id}/edit", postEditCard.Handle())
	http.HandleFunc("DELETE /sheet/{section}/{card_id}", deleteCard.Handle())
	http.HandleFunc("POST /sheet/{section}/add", addCardHandler.Handle())
	http.HandleFunc("GET /sheet/{section}", sectionViewHandler.Handle())

	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
