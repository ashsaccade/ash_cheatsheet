package main

import (
	"html/template"
	"net/http"
)

type SectionData struct {
	Name string
}

func main() {
	tmpl, err := template.ParseGlob("templates/*.htm")
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

			if err := sectionHtm.Execute(writer, SectionData{
				Name: section,
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
