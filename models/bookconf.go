package models

type BookConf struct {
	ID      string
	Url     string
	Format  string
	Path    string
	Jump    int
	Delay   int
	IsDesc  bool
	Catalog string
	// Font     string
}
