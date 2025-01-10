package getsection

import (
	"ash_cheatsheet/internal/cards"
	"ash_cheatsheet/internal/entities"
	"ash_cheatsheet/internal/render"
	"github.com/pav5000/go-common/lambda"
	"github.com/pav5000/logger"
	"html/template"
	"net/http"
)

type Handler struct {
	cards       *cards.Service
	htmTemplate *template.Template
}

type SectionData struct {
	SectionName string
	Cards       []CardView
}

type CardView struct {
	Name        string
	Description template.HTML
	Id          int64
}

func New(tmpl *template.Template, c *cards.Service) *Handler {
	sectionHtm := tmpl.Lookup("section.htm")
	if sectionHtm == nil {
		panic("section htm is nil")
	}

	return &Handler{cards: c, htmTemplate: sectionHtm}
}

func (h *Handler) Handle() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		section := request.PathValue("section")

		cardsBySection, err := h.cards.GetCards(section)
		if err != nil {
			logger.Fatal(err.Error())
		}

		err = h.htmTemplate.Execute(writer, SectionData{
			SectionName: section,
			Cards: lambda.MapSlice(cardsBySection, func(card entities.Card) CardView {
				return CardView{
					Name:        card.Name,
					Description: template.HTML(render.Render(card.Description)),
					Id:          card.Id,
				}
			}),
		})
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}
