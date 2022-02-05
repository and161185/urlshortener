package usstorage

import (
	//"os"

	"os"
	"testing"
	"urlshortener/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGenerateShortUrl(t *testing.T) {
	dbname := "test_gsu.db"
	log := getLog()
	os.Remove("./database/" + dbname)
	d := NewUSStorage(log, dbname)
	defer os.Remove("./database/" + dbname)

	us := models.FullUrlScheme{
		Url: "http:\\yandex.ru",
	}
	res, _ := d.GenerateShortUrl(us)
	assert.Equal(t, "AQ", res.ShortId)

}

func TestGetFullUrl(t *testing.T) {
	dbname := "test_gfu.db"
	log := getLog()
	os.Remove("./database/" + dbname)
	d := NewUSStorage(log, dbname)
	os.Remove("./database/" + dbname)

	us := models.FullUrlScheme{
		Url: "http:\\yandex.ru",
	}
	su, _ := d.GenerateShortUrl(us)
	res, _ := d.GetFullUrl(su.ShortId)

	assert.Equal(t, "http:\\yandex.ru", res.Url)
}

func TestGetStats(t *testing.T) {
	dbname := "test_gs.db"
	log := getLog()
	os.Remove("./database/" + dbname)
	d := NewUSStorage(log, dbname)
	os.Remove("./database/" + dbname)

	us := models.FullUrlScheme{
		Url: "http:\\yandex.ru",
	}
	su, _ := d.GenerateShortUrl(us)
	d.RegisterClick(su.ShortId, "127.0.0.1")
	stats, _ := d.GetStats(su.StatId)

	assert.Equal(t, int64(1), stats.ClickCount)
	assert.Equal(t, "127.0.0.1", stats.Clicks[0].IP)

}

func getLog() *logrus.Logger {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	Formatter := new(logrus.TextFormatter)

	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	log.Formatter = new(logrus.JSONFormatter)

	return log
}
