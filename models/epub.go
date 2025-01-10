package models

type Section struct {
	Url     string
	Title   string
	Content string
	IsVol   bool
}

type Assets struct {
	MainCss        string
	CoverCss       string
	FontCss        string
	InstructionCss string
	Font           string
	MetaInf        string
}

type Images struct {
	Cover       string
	ColImg      string
	IntroImg    string
	ContentLogo string
	VolImg      string
}
