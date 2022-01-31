package usrepo

import (
	"fmt"
	"urlshortener/internal/models"
)

type UrlShortenerRepo interface {
	Close()
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

func (us *UrlShortener) Close() {
	us.repo.Close()
}

func (us *UrlShortener) GenerateShortUrl(url models.FullUrlScheme) (data *models.ShortLinkScheme, err error) {
	data, err = us.repo.GenerateShortUrl(url)
	if err != nil {
		return nil, fmt.Errorf("generate short url error: %w", err)
	}

	return data, nil
}

func (us *UrlShortener) GetFullUrl(shortId string) (urlScheme *models.FullUrlScheme, err error) {
	urlScheme, err = us.repo.GetFullUrl(shortId)
	if err != nil {
		return nil, fmt.Errorf("get full url error: %w", err)
	}

	return urlScheme, nil
}

func (us *UrlShortener) RegisterClick(shortId string, ip string) (err error) {
	err = us.repo.RegisterClick(shortId, ip)
	if err != nil {
		return fmt.Errorf("register click error: %w", err)
	}

	return nil
}

func (us *UrlShortener) GetStats(statId string) (ss *models.StatsScheme, err error) {
	ss, err = us.repo.GetStats(statId)
	if err != nil {
		return nil, fmt.Errorf("get stats error: %w", err)
	}

	return ss, nil
}
