package cards

import (
	"ash_cheatsheet/internal/entities"
	"ash_cheatsheet/internal/errs"
	"context"
	"strings"
)

type Repository interface {
	InsertNewCard(ctx context.Context, card entities.Card) error
	DeleteCard(ctx context.Context, id int64, sectionName string) error
	GetCardByID(ctx context.Context, id int64) (*entities.Card, error)
	UpdateCard(ctx context.Context, id int64, name, description string) error
}

type Client struct {
	db Repository
}

func New(db Repository) *Client { return &Client{db: db} }

func (c *Client) DeleteCard(id int64, sectionName string) error {
	if err := c.db.DeleteCard(context.Background(), id, sectionName); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateCardByID(id int64, name, description string) error {
	err := c.db.UpdateCard(context.Background(), id, name, description)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetCardByID(id int64) (*entities.Card, error) {
	card, err := c.db.GetCardByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return card, err
}

func (c *Client) CreateNewCard(name, description, sectionName string) error {
	if strings.TrimSpace(name) == "" {
		return errs.ErrEmptyCardName
	}

	if strings.TrimSpace(description) == "" {
		return errs.ErrEmptyCardDesc
	}

	err := c.db.InsertNewCard(context.Background(), entities.Card{
		Name:        name,
		Description: description,
		Section:     sectionName,
	})
	if err != nil {
		panic(err)
	}
	return err
}
