package models

type Page struct {
	Title string
	Dom   string
}

type Style struct {
	Instruction Page
	Desc        Page
	Vol         string
	Chapter     string
}

type UrlWithHeader struct {
	Url        string
	Header     map[string]string
	NeedDivide bool `toml:"need_divide"`
}

type BookFilter struct {
	Selector string
	Filter   []string
}

type ChapterFilter struct {
	Element string
	Url     string
	Title   string
}

type BookCatch struct {
	Domain    string
	DelayTime int `toml:"delay_time"`
	Cover     UrlWithHeader
	Name      BookFilter
	ID        string
	Toc       string
	Sort      string
	Title     BookFilter
	Author    BookFilter
	Intro     BookFilter
	Chapter   ChapterFilter
	Content   BookFilter
}

type InfoSelector struct {
	Catalog     string   `toml:"catalog"`
	VolName     string   `toml:"vol_name"`
	Chapter     string   `toml:"chapter"`
	Api         string   `toml:"api"`
	IsJSON      bool     `toml:"is_json"`
	PassVols    []string `toml:"pass_vols"`
	ExcludeVols []string `toml:"exclude_vols"`
}

type Config struct {
	*Style
	BookCatch    map[string]*BookCatch    `toml:"book_catch"`
	InfoSelector map[string]*InfoSelector `toml:"info_selector"`
	TmpDir       string                   `toml:"-"`
	DelayTime    int                      `toml:"delay_time"`
	From         string
}
