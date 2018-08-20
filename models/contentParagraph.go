package models

import "github.com/go-ozzo/ozzo-validation"

type ParagraphType int

var ParagraphTypes = struct{
	PlainText ParagraphType
	Hyperlink ParagraphType
}{
	PlainText:1,
	Hyperlink:2,
}

type Paragraph struct {
	Id      int           `json:"id" db:"id"`
	SeedId  int           `json:"seedId" db:"seed_id"`
	Type    ParagraphType `json:"type" db:"type"`
	Content string        `json:"content" db:"content"`
}

func (m Paragraph) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Content, validation.Required),
		validation.Field(&m.SeedId, validation.Required),
	)
}
