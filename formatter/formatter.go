package formatter

type Formatter interface {
	InitBook() error
	GenContentPrefix(int, string)
	GenBookContent(int) error
	Build() error
}
