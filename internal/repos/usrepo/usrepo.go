package usrepo

import (
	"fmt"
	"urlshortener/internal/models"
)

type UrlShortenerRepo interface {
	GenerateShortUrl(url models.FullUrlScheme) (data *models.ShortLinkScheme, err error)
	GetFullUrl(shortId string) (urlScheme *models.FullUrlScheme, err error)
	RegisterClick(shortId string, ip string) (err error)
	GetStats(statId string) (ss *models.StatsScheme, err error)
}

type UrlShortener struct {
	repo UrlShortenerRepo
}

func NewUrlShortener(r UrlShortenerRepo) *UrlShortener {
	return &UrlShortener{
		repo: r,
	}
}

//GenerateShortUrl returns scheme with shortId and relative data
func (us *UrlShortener) GenerateShortUrl(url models.FullUrlScheme) (data *models.ShortLinkScheme, err error) {
	data, err = us.repo.GenerateShortUrl(url)
	if err != nil {
		return nil, fmt.Errorf("generate short url error: %w", err)
	}

	return data, nil
}

//GetFullUrl converts short id into full url for redirect
func (us *UrlShortener) GetFullUrl(shortId string) (urlScheme *models.FullUrlScheme, err error) {
	urlScheme, err = us.repo.GetFullUrl(shortId)
	if err != nil {
		return nil, fmt.Errorf("get full url error: %w", err)
	}

	return urlScheme, nil
}

//RegisterClick collects statistics for shortId
func (us *UrlShortener) RegisterClick(shortId string, ip string) (err error) {
	err = us.repo.RegisterClick(shortId, ip)
	if err != nil {
		return fmt.Errorf("register click error: %w", err)
	}

	return nil
}

//RegisterClick statistics scheme for shortId using statId
func (us *UrlShortener) GetStats(statId string) (ss *models.StatsScheme, err error) {
	ss, err = us.repo.GetStats(statId)
	if err != nil {
		return nil, fmt.Errorf("get stats error: %w", err)
	}

	return ss, nil
}
