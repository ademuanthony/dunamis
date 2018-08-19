package models

import (
	"github.com/go-ozzo/ozzo-validation"
)

//Seed represent a day Seed of destiny record
type Seed struct {
	Id           int         `json:"id" db:"id"`
	Title        string      `json:"title" db:"title"`
	Thought      string      `json:"thought" db:"thought"`
	Content      string      `json:"content" db:"content"`
	Prayer       string      `json:"prayer" db:"prayer"`
	Assignment   string      `json:"assignment" db:"assignment"`
	DailyReading string      `json:"daily_reading" db:"daily_reading"`
	Quote        string      `json:"quote" db:"quote"`
	Resource     string      `json:"resource" db:"resource"`
	Date         string      `json:"date" db:"date"`
	Year         int         `json:"year" db:"year"`
	Month        int         `json:"month" db:"month"`
	Day          int         `json:"day" db:"day"`
	Paragraphs   []Paragraph `json:"paragraphs" db:"-"`
	Scripture    string      `json:"scripture" db:"scripture"`
	RememberThis string      `json:"remember_this" db:"remember_this"`
}

func (m Seed) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required),
		//validation.Field(&m.Content, validation.Required),
		validation.Field(&m.Date, validation.Required),
	)
}