package source

import "freb/models"

type Source interface {
	GetBook(*models.Book) error
}
