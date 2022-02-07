package usrepo

import (
	"testing"
	"urlshortener/internal/models"

	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
}

func (m *mockStorage) GenerateShortUrl(url models.FullUrlScheme) (data *models.ShortLinkScheme, err error) {
	return &models.ShortLinkScheme{ShortId: "AQ"}, nil
}

func (m *mockStorage) RegisterClick(shortId string, ip string) (err error) {
	return nil
}

func (m *mockStorage) GetFullUrl(shortId string) (urlScheme *models.FullUrlScheme, err error) {
	return &models.FullUrlScheme{Url: "http:\\yandex.ru"}, nil
}

func (m *mockStorage) GetStats(statId string) (ss *models.StatsScheme, err error) {
	var clicks []*models.ClickScheme
	click := &models.ClickScheme{
		IP: "127.0.0.1",
	}
	clicks = append(clicks, click)

	return &models.StatsScheme{ClickCount: int64(1), Clicks: clicks}, nil
}

func TestGenerateShortUrl(t *testing.T) {

	d := &mockStorage{}
	us := NewUrlShortener(d)

	fus := models.FullUrlScheme{
		Url: "http:\\yandex.ru",
	}
	res, _ := us.GenerateShortUrl(fus)
	assert.Equal(t, "AQ", res.ShortId)
}

func TestGetFullUrl(t *testing.T) {
	d := &mockStorage{}
	us := NewUrlShortener(d)

	fus := models.FullUrlScheme{
		Url: "http:\\yandex.ru",
	}
	su, _ := us.GenerateShortUrl(fus)
	res, _ := us.GetFullUrl(su.ShortId)

	assert.Equal(t, "http:\\yandex.ru", res.Url)

}

func TestGetStats(t *testing.T) {
	d := &mockStorage{}
	us := NewUrlShortener(d)

	fus := models.FullUrlScheme{
		Url: "http:\\yandex.ru",
	}
	su, _ := us.GenerateShortUrl(fus)
	us.RegisterClick(su.ShortId, "127.0.0.1")
	stats, _ := d.GetStats(su.StatId)

	assert.Equal(t, int64(1), stats.ClickCount)
	assert.Equal(t, "127.0.0.1", stats.Clicks[0].IP)
}
