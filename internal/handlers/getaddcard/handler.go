package getaddcard

import (
	"net/http"

	"html/template"

	"github.com/pav5000/logger"
)

type AddCardData struct {
	SectionName string
}

type Handler struct {
	htmTemplate *template.Template
}

func New(tmpl *template.Template) *Handler {
	cardHtm := tmpl.Lookup("add_card.htm")
	if cardHtm == nil {
		panic("new_card htm is nil")
	}

	return &Handler{htmTemplate: cardHtm}
}

func (h Handler) Handle() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		section := request.PathValue("section")

		if err := h.htmTemplate.Execute(writer, AddCardData{SectionName: section}); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
