package add

import (
	"net/http"

	"html/template"

	"github.com/pav5000/logger"
)

type AddCardData struct {
	SectionName string
}

func New(tmpl template.Template) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		cardHtm := tmpl.Lookup("add_card.htm")
		if cardHtm == nil {
			panic("new_card htm is nil")
		}
		section := request.PathValue("section")

		if err := cardHtm.Execute(writer, AddCardData{SectionName: section}); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
