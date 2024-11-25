package static

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pav5000/logger"
)

func New() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fileName := request.PathValue("filename")

		switch {
		case strings.HasSuffix(fileName, ".css"):
			writer.Header().Add("Content-Type", "text/css")
		case strings.HasSuffix(fileName, ".woff2"):
			writer.Header().Add("Content-Type", "application/font-woff2")
		case strings.HasSuffix(fileName, ".js"):
			writer.Header().Add("Content-Type", "text/javascript")
		default:
			panic("invalid file")
		}

		f, err := os.Open("static/" + fileName)
		if err != nil {
			logger.Fatal(err.Error())
		}
		defer f.Close()

		if _, err := io.Copy(writer, f); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
