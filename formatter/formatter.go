package formatter

import (
	"bytes"
)

type Formatter interface {
	InitBook() error
	GenContentPrefix(*bytes.Buffer, string)
	GenBookContent(int, bytes.Buffer, bytes.Buffer) error
	Build() error
}
