package db

import (
	"ash_cheatsheet/internal/entities"
	"context"
	"embed"
	"time"

	"github.com/pav5000/easy-sqlite"
	"github.com/pav5000/go-common/lambda"
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
	UpdatedAt   time.Time `db:"updated_at"`
	Pinned      bool
}

func (c *Client) UpdateCard(ctx context.Context, id int64, name, description string) error {
	q := `update cards set name=?, description=? where id = ?`

	_, err := c.conn.ExecContext(ctx, q, name, description, id)
	return err
}

func (c *Client) TogglePinCard(ctx context.Context, id int64, sectionName string) error {
	q := `update cards set pinned = not pinned where id = ? and section = ?`

	_, err := c.conn.ExecContext(ctx, q, id, sectionName)
	return err
}

func (c *Client) GetCardByID(ctx context.Context, id int64) (*entities.Card, error) {
	q := `select id, name, description, section, pinned from cards where id = ?`

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
		Pinned:      dbCard.Pinned,
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
	insert into cards (name, description, section, updated_at, pinned)
	values (?, ?, ?, ?, ?)`

	_, err := c.conn.ExecContext(ctx, q, card.Name, card.Description, card.Section, card.UpdatedAt, card.Pinned)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SelectAllCardsBySection(ctx context.Context, sectionName string) ([]entities.Card, error) {
	q := `
		select id, name, description, section, updated_at, pinned
		from cards
		where section = ? order by pinned desc, updated_at desc`

	var cardsDb []CardRow

	err := c.conn.SelectContext(ctx, &cardsDb, q, sectionName)
	if err != nil {
		return nil, err
	}

	res := lambda.MapSlice(cardsDb, func(row CardRow) entities.Card {
		return entities.Card{
			Id:          row.Id,
			Name:        row.Name,
			Description: row.Description,
			Section:     row.Section,
			UpdatedAt:   row.UpdatedAt,
			Pinned:      row.Pinned,
		}
	})
	return res, err
}
