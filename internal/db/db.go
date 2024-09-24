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

func (c *Client) UpdateCard(ctx context.Context, id int64, name, description string) error {
	q := `update cards set name=?, description=? where id = ?`

	_, err := c.conn.ExecContext(ctx, q, name, description, id)
	return err
}

func (c *Client) GetCardByID(ctx context.Context, id int64) (*entities.Card, error) {
	q := `select id, name, description, section from cards where id = ?`

	dbCard := CardRow{}
	err := c.conn.GetContext(ctx, &dbCard, q, id)
	if err != nil {
		return nil, err
	}

	return &entities.Card{
		Id:          dbCard.Id,
		Name:        dbCard.Name,
		Description: dbCard.Description,
		Section:     dbCard.Section,
	}, nil
}

func (c *Client) DeleteCard(ctx context.Context, id int64, sectionName string) error {
	q := `delete from cards where id = ? and section = ?`

	_, err := c.conn.ExecContext(ctx, q, id, sectionName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) InsertNewCard(ctx context.Context, card entities.Card) error {
	q := `
	insert into cards (name, description, section)
	values (?, ?, ?)`

	_, err := c.conn.ExecContext(ctx, q, card.Name, card.Description, card.Section)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SelectAllCardsBySection(ctx context.Context, sectionName string) ([]entities.Card, error) {
	q := `
		select id, name, description, section
		from cards
		where section = ?`

	var cardsDb []CardRow

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
