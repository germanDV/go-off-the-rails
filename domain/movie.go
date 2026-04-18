package domain

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Movie struct {
	ID        uuid.UUID
	OrgID     uuid.UUID
	Title     string
	Rating    int
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMovie(orgID uuid.UUID, titleInput string, ratingInput string) (Movie, error) {
	title := strings.TrimSpace(titleInput)

	rating, err := strconv.Atoi(strings.TrimSpace(ratingInput))
	if err != nil {
		return Movie{}, err
	}

	err = validate(title, rating)
	if err != nil {
		return Movie{}, err
	}

	return Movie{
		ID:        uuid.Must(uuid.NewV7()),
		OrgID:     orgID,
		Title:     title,
		Rating:    rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (m *Movie) Update(titleInput string, ratingInput string) error {
	title := strings.TrimSpace(titleInput)

	rating, err := strconv.Atoi(strings.TrimSpace(ratingInput))
	if err != nil {
		return err
	}

	err = validate(title, rating)
	if err != nil {
		return err
	}

	m.Title = title
	m.Rating = rating
	m.UpdatedAt = time.Now().UTC()

	return nil
}

func validate(title string, rating int) error {
	if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}

	if rating < 1 || rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}

	return nil
}
