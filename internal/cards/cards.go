package cards

import (
	"ash_cheatsheet/internal/entities"
	"ash_cheatsheet/internal/errs"
	"context"
	"fmt"
	"strings"
	"time"
)

type Repository interface {
	InsertNewCard(ctx context.Context, card entities.Card) error
	DeleteCard(ctx context.Context, id int64, sectionName string) error
	GetCardByID(ctx context.Context, id int64) (*entities.Card, error)
	UpdateCard(ctx context.Context, id int64, name, description string) error
	SelectAllCardsBySection(ctx context.Context, sectionName string) ([]entities.Card, error)
}

type Service struct {
	db Repository
}

func New(db Repository) *Service { return &Service{db: db} }

func (s *Service) DeleteCard(id int64, sectionName string) error {
	return s.db.DeleteCard(context.Background(), id, sectionName)
}

func (s *Service) UpdateCardByID(id int64, name, description string) error {
	return s.db.UpdateCard(context.Background(), id, name, description)
}

func (s *Service) GetCardByID(id int64) (*entities.Card, error) {
	return s.db.GetCardByID(context.Background(), id)
}

func (s *Service) CreateNewCard(name, description, sectionName string) error {
	if strings.TrimSpace(name) == "" {
		return errs.ErrEmptyCardName
	}

	if strings.TrimSpace(description) == "" {
		return errs.ErrEmptyCardDesc
	}

	err := s.db.InsertNewCard(context.Background(), entities.Card{
		Name:        name,
		Description: description,
		Section:     sectionName,
		LastUpdated: time.Now(),
	})
	return err
}

func (s *Service) GetCards(sectionName string) ([]entities.Card, error) {
	res, err := s.db.SelectAllCardsBySection(context.Background(), sectionName)
	if err != nil {
		return nil, fmt.Errorf("db.SelectAllCardsBySection: %w", err)
	}
	return res, nil
}
