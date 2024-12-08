package posteditcard

import (
	"ash_cheatsheet/internal/cards"
	"github.com/pav5000/logger"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type EditCardResult struct {
	Message string
	Class   string
}

type Handler struct {
	cards       *cards.Service
	htmTemplate *template.Template
}

func New(tmpl *template.Template, cards *cards.Service) *Handler {
	editCardResultHtm := tmpl.Lookup("edit_result.htm")
	if editCardResultHtm == nil {
		panic("edit_result htm is nil")
	}

	return &Handler{cards: cards, htmTemplate: editCardResultHtm}
}

func (h Handler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := h.cards.UpdateCardByID(int64(cardId), name, description); err != nil {
			logger.Fatal(err.Error())
		}

		if err := h.htmTemplate.Execute(w, EditCardResult{Message: "Карточка изменена"}); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
