package deletecard

import (
	"ash_cheatsheet/internal/cards"
	"net/http"
	"strconv"

	"github.com/pav5000/logger"
)

type Handler struct {
	cards *cards.Service
}

func New(c *cards.Service) Handler { return Handler{cards: c} }

func (h Handler) Handle() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		sectionName := request.PathValue("section")
		idStr := request.PathValue("card_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if err := h.cards.DeleteCard(int64(id), sectionName); err != nil {
			logger.Fatal(err.Error())
		}
	}
}
