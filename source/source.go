package source

import (
	"freb/formatter"
)

type Source interface {
	GetBook(formatter.EpubFormat) error
}
