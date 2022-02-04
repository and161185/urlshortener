package models

type ShortLinkScheme struct {
	FullUrl        string
	ShortId        string
	StatId         string
	ExpirationDate string
}

type StatsScheme struct {
	ClickCount     int64
	ExpirationDate string
	Clicks         []*ClickScheme
}

type FullUrlScheme struct {
	Url string
}

type ClickScheme struct {
	IP   string
	Time string
}
