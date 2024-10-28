package preview

import (
	"ash_cheatsheet/internal/render"
	"io"
	"net/http"
	"net/url"

	"github.com/pav5000/logger"
)

func New() func(writer http.ResponseWriter, r *http.Request) {
	return func(writer http.ResponseWriter, r *http.Request) {
		rawFormData, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Fatal(err.Error())
		}

		vals, err := url.ParseQuery(string(rawFormData))
		if err != nil {
			logger.Fatal(err.Error())
		}

		description := vals.Get("description")

		res := render.Render(description)
		writer.Write([]byte(res))
	}
}
