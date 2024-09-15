package db

import (
	"ash_cheatsheet/internal/entities"
	"context"
	"embed"
	"github.com/pav5000/easy-sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Client struct {
	conn *easysqlite.EasySqlite
}

func New() (*Client, error) {
	db, err := easysqlite.New("data/db.sqlite", embedMigrations, "migrations")
	if err != nil {
		return nil, err
	}

	return &Client{conn: db}, nil
}

type CardRow struct {
	Id          int64
	Name        string
	Description string
	Section     string
}

func (c *Client) SelectAllCardsBySection(ctx context.Context, sectionName string) ([]entities.Card, error) {
	q := `
		select id, name, description, section
		from cards
		where section = ?`

	cardsDb := []CardRow{}

	err := c.conn.SelectContext(ctx, &cardsDb, q, sectionName)
	if err != nil {
		return nil, err
	}

	res := make([]entities.Card, 0, len(cardsDb))
	for _, row := range cardsDb {
		res = append(res, entities.Card{
			Id:          row.Id,
			Name:        row.Name,
			Description: row.Description,
			Section:     row.Section,
		})
	}
	return res, nil
}
