package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Movie is a database schema.
type Movie struct {
	ID          uint      `gorm:"column:id"           json:"id"`
	Created     time.Time `gorm:"column:created"      json:"created"`
	Modified    time.Time `gorm:"column:modified"     json:"modified"`
	Title       string    `gorm:"column:title"        json:"title"`
	Description string    `gorm:"column:description"  json:"description"`
	ReleaseYear uint      `gorm:"column:release_year" json:"release_year"`
	IMDbRating  float64   `gorm:"column:imdb_rating"  json:"imdb_rating"`
	Link        string    `gorm:"column:link"         json:"link"`
}

// ListMovieByYear returns movies from a given year.
func ListMovieByYear(year string) ([]*Movie, error) {
	result := []*Movie{}

	if err := PSQL.Where("release_year = ?", year).Order("imdb_rating DESC").Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// TableName explicitly tells the ORM where to look for the table for this row.
func (Movie) TableName() string {
	return "movies"
}

// BeforeUpdate is a hook to database transaction for setting timestamp.
func (m *Movie) BeforeUpdate(scope *gorm.Scope) {
	scope.SetColumn("Modified", time.Now().Round(time.Microsecond))
}

// BeforeCreate is a hook to database transaction for setting timestamp.
func (m *Movie) BeforeCreate(scope *gorm.Scope) {
	scope.SetColumn("Created", time.Now().Round(time.Microsecond))
	scope.SetColumn("Modified", time.Now().Round(time.Microsecond))
}

// Create inserts a movie to database.
func (m *Movie) Create() error {
	if PSQL == nil {
		return errors.New("cannot create without initializing database connection")
	}

	return PSQL.Create(m).Error
}
