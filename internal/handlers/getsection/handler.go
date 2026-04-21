package getsection

import (
	"ash_cheatsheet/internal/cards"
	"ash_cheatsheet/internal/entities"
	"ash_cheatsheet/internal/render"
	"html/template"
	"net/http"
	"time"

	"github.com/pav5000/go-common/lambda"
	"github.com/pav5000/logger"
)

type Handler struct {
	cards       *cards.Service
	htmTemplate *template.Template
}

type SectionData struct {
	SectionName string
	Sections    []SectionView
	Cards       []CardView
}

type SectionView struct {
	Name   string
	Active bool
}

type CardView struct {
	Name        string
	Description template.HTML
	Id          int64
	UpdatedAt   string
	Pinned      bool
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

		sections, err := h.cards.GetSections()
		if err != nil {
			logger.Fatal(err.Error())
		}

		err = h.htmTemplate.Execute(writer, SectionData{
			SectionName: section,
			Sections: lambda.MapSlice(sections, func(sectionName string) SectionView {
				return SectionView{
					Name:   sectionName,
					Active: sectionName == section,
				}
			}),
			Cards: lambda.MapSlice(cardsBySection, func(card entities.Card) CardView {
				v := CardView{
					Name:      card.Name,
					Id:        card.Id,
					UpdatedAt: card.UpdatedAt.Format(time.DateOnly),
					Pinned:    card.Pinned,
				}
				renderedDesc, err := render.Render(card.Description)
				if err != nil {
					logger.Error(err.Error())
					return v
				}
				v.Description = template.HTML(renderedDesc)
				return v
			}),
		})
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}
