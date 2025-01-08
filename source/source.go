package source

import (
	"freb/formatter"
	"freb/models"
)

type Source interface {
	GetBook(*formatter.EpubFormat, *models.BookCatch) error
}
