package conf

import "regexp"

var (
	Conf *Config
)

var PrivacyReg []*regexp.Regexp

var (
	// StoragesLoaded loaded success if empty
	StoragesLoaded = false
)
var (
	RawIndexHtml string
	ManageHtml   string
	IndexHtml    string
)
