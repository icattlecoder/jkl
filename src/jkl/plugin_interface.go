package jkl

type plugin_interface interface {
	PlugeName() string
	SetArgs(args []string)
	OnPageParsed(page *Page)
	OnSiteGenarated(site *Site)
}
