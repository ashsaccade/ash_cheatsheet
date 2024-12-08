package addcard

import (
	"ash_cheatsheet/internal/cards"
	"ash_cheatsheet/internal/errs"
	"github.com/pav5000/logger"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type Handler struct {
	cards *cards.Service
	tmpl  *template.Template
}

type AddCardResult struct {
	Message string
	Class   string // is-danger
}

func New(c *cards.Service, tmpl *template.Template) *Handler {
	addResult := tmpl.Lookup("add_result.htm")
	if addResult == nil {
		panic("add result htm is nil")
	}

	return &Handler{cards: c, tmpl: addResult}
}

func (h Handler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
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

		err = h.cards.CreateNewCard(vals.Get("name"), vals.Get("description"), sectionName)
		switch {
		case errors.Is(err, errs.ErrEmptyCardName):
			res = AddCardResult{Message: "Имя карточки не может быть пустым!", Class: "is-danger"}
		case errors.Is(err, errs.ErrEmptyCardDesc):
			res = AddCardResult{Message: "Описание карточки не может быть пустым!", Class: "is-danger"}
		default:
			res = AddCardResult{Message: "Новая карточка добавлена в секцию " + sectionName}
		}
		if err := h.tmpl.Execute(writer, res); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
