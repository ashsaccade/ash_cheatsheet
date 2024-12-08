package geteditcard

import (
	"ash_cheatsheet/internal/cards"
	"github.com/pav5000/logger"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	client      *cards.Service
	htmTemplate *template.Template
}

type EditCardPage struct {
	Id          int64
	Name        string
	Description string
	SectionName string
}

func New(tmpl *template.Template, client *cards.Service) Handler {
	editCardHtm := tmpl.Lookup("edit.htm")
	if editCardHtm == nil {
		panic("edit card htm is nil")
	}

	return Handler{client, editCardHtm}
}

func (h Handler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("card_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Fatal(err.Error())
		}

		card, err := h.client.GetCardByID(int64(id))
		if err != nil {
			logger.Fatal(err.Error())
		}

		if err := h.htmTemplate.Execute(w, EditCardPage{
			Id:          card.Id,
			Name:        card.Name,
			Description: card.Description,
			SectionName: card.Section,
		}); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
