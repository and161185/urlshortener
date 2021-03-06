package usstorage

import (
	"database/sql"
	"errors"
	"math/big"
	neturl "net/url"
	"strings"
	"time"

	b64 "encoding/base64"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"urlshortener/internal/models"
)

//CreateUrlsTable creates urls table (if doesn't exists) with fields:
//id INTEGER, shortId TEXT, statId TEXT, url TEXT, expirationDate TIME
func CreateUrlsTableSqlite3(db *sql.DB, log *logrus.Logger) {

	log.Info("Creating urls table")

	checkTableSQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='urls';"
	row, err := db.Query(checkTableSQL)
	if err != nil {
		log.Fatal("Checking if urls table exists ", err)
	}
	defer row.Close()

	urlsTabelExists := false
	for row.Next() {
		urlsTabelExists = true
	}

	if !urlsTabelExists {
		createStudentTableSQL := `CREATE TABLE urls (
			id		INTEGER PRIMARY KEY AUTOINCREMENT
								UNIQUE
								NOT NULL,
			shortId	TEXT    NOT NULL
								UNIQUE,
			statId	TEXT    NOT NULL
								UNIQUE,
			url		TEXT    NOT NULL,
			expirationDate TIME
		);`

		log.Info("Create urls table...")
		statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
		if err != nil {
			log.Fatal("Creating urls table ", err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal("can't create urls table", err)
		}
		log.Info("urls table created")
	} else {
		log.Info("urls table exists")
	}
}

//CreateClicksTable creates clicks table (if doesn't exists) with fields:
//shortId TEXT, IP TEXT, time TIME
//clicks table collects stats for shortId
func CreateClicksTableSqlite3(db *sql.DB, log *logrus.Logger) {
	checkTableSQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='clicks';"
	row, err := db.Query(checkTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	urlsTabelExists := false
	for row.Next() {
		urlsTabelExists = true
	}

	if !urlsTabelExists {
		createClicksTableSQL := `
		CREATE TABLE clicks (
			shortId TEXT NOT NULL
						 REFERENCES urls (shortId) ON DELETE CASCADE,
			IP      TEXT NOT NULL,
			time    TIME NOT NULL
		);`

		log.Info("Create clicks table...")
		statement, err := db.Prepare(createClicksTableSQL) // Prepare SQL Statement
		if err != nil {
			log.Fatal(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal("can't create clicks table", err)
		}
		log.Info("clicks table created")
	} else {
		log.Info("clicks table exists")
	}
}

//CreateUrlsTable creates urls table (if doesn't exists) with fields:
//id INTEGER, shortId TEXT, statId TEXT, url TEXT, expirationDate TIME
func CreateUrlsTablePostgres(db *sql.DB, log *logrus.Logger) {

	log.Info("Creating urls table")

	checkTableSQL := "SELECT tablename FROM pg_tables WHERE tablename  = 'urls';"
	row, err := db.Query(checkTableSQL)
	if err != nil {
		log.Fatal("Checking if urls table exists ", err)
	}
	defer row.Close()

	urlsTabelExists := false
	for row.Next() {
		urlsTabelExists = true
	}

	if !urlsTabelExists {
		createStudentTableSQL := `CREATE TABLE urls (
			id		INTEGER PRIMARY KEY AUTOINCREMENT
								UNIQUE
								NOT NULL,
			shortId	TEXT    NOT NULL
								UNIQUE,
			statId	TEXT    NOT NULL
								UNIQUE,
			url		TEXT    NOT NULL,
			expirationDate TIME
		);`

		log.Info("Create urls table...")
		statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
		if err != nil {
			log.Fatal("Creating urls table ", err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal("can't create urls table", err)
		}
		log.Info("urls table created")
	} else {
		log.Info("urls table exists")
	}
}

//CreateClicksTable creates clicks table (if doesn't exists) with fields:
//shortId TEXT, IP TEXT, time TIME
//clicks table collects stats for shortId
func CreateClicksTablePostgres(db *sql.DB, log *logrus.Logger) {
	checkTableSQL := "SELECT tablename FROM pg_tables WHERE tablename  = 'clicks';"
	row, err := db.Query(checkTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	urlsTabelExists := false
	for row.Next() {
		urlsTabelExists = true
	}

	if !urlsTabelExists {
		createClicksTableSQL := `
		CREATE TABLE clicks (
			shortId TEXT NOT NULL
						 REFERENCES urls (shortId) ON DELETE CASCADE,
			IP      TEXT NOT NULL,
			time    TIME NOT NULL
		);`

		log.Info("Create clicks table...")
		statement, err := db.Prepare(createClicksTableSQL) // Prepare SQL Statement
		if err != nil {
			log.Fatal(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal("can't create clicks table", err)
		}
		log.Info("clicks table created")
	} else {
		log.Info("clicks table exists")
	}
}

//Close calls Close method of *sql.DB
func (d *dbdriver) Close() {
	d.db.Close()
}

//GenerateShortUrl inserts new row into urls table
//returns scheme with shortId and relative data
func (d *dbdriver) GenerateShortUrl(url models.FullUrlScheme) (data *models.ShortLinkScheme, err error) {

	_, err = neturl.ParseRequestURI(url.Url)
	if err != nil {
		d.log.Error(err)
		return nil, err
	}

	statId := NewStatKey()
	shortId := statId

	d.log.Info("Inserting url record ", statId)

	tx, err := d.db.Begin()
	if err != nil {
		d.log.Error(err)
		return nil, err
	}

	expirationDate := time.Now().AddDate(0, 1, 0)

	insertSQL := `INSERT INTO urls(statId, shortId, url, expirationDate) VALUES (?, ?, ?, ?)`
	sqlResult, err := tx.Exec(insertSQL, statId, shortId, url.Url, expirationDate)
	if err != nil {
		d.log.Error(err)
		err = tx.Rollback()
		if err != nil {
			d.log.Fatal(err)
		}
		return nil, err
	}

	LastInsertedId, err := sqlResult.LastInsertId()
	if err != nil {
		d.log.Error(err)
		err = tx.Rollback()
		if err != nil {
			d.log.Fatal(err)
		}
		return nil, err
	}

	shortId = getShortId(LastInsertedId)
	updateSql := `UPDATE urls SET shortId = ? WHERE id = ?`
	_, err = tx.Exec(updateSql, shortId, LastInsertedId)
	if err != nil {
		d.log.Error(err)
		err = tx.Rollback()
		if err != nil {
			d.log.Fatal(err)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		d.log.Error(err)
		err = tx.Rollback()
		if err != nil {
			d.log.Fatal(err)
		}
		return nil, err
	}

	result := &models.ShortLinkScheme{
		FullUrl:        url.Url,
		ShortId:        shortId,
		StatId:         statId,
		ExpirationDate: expirationDate.Format("2006-01-02"),
	}

	return result, nil
}

//GetFullUrl converts short id into full url
func (d *dbdriver) GetFullUrl(shortId string) (urlScheme *models.FullUrlScheme, err error) {
	query := `select url from urls WHERE shortId = ?`
	rows := d.db.QueryRow(query, shortId)

	var fullUrl string
	err = rows.Scan(&fullUrl)

	if err == sql.ErrNoRows {
		error := errors.New("short url doesn't exist")
		d.log.Error(error)
		return nil, error
	} else if err != nil {
		d.log.Error(err)
		return nil, err
	}

	urlScheme = &models.FullUrlScheme{Url: fullUrl}

	return urlScheme, nil
}

//RegisterClick inserts new row into clicks table
func (d *dbdriver) RegisterClick(shortId string, ip string) (err error) {
	insertSQL := `INSERT INTO clicks(shortId, IP, time) VALUES (?, ?, ?)`
	_, err = d.db.Exec(insertSQL, shortId, ip, time.Now())

	if err != nil {
		d.log.Error(err)
	}

	return err
}

//GetStats return stats scheme for short link using statId
func (d *dbdriver) GetStats(statId string) (ss *models.StatsScheme, err error) {
	query := `SELECT urls.ShortID, MAX(urls.expirationDate) as expirationDate, COALESCE(count(clicks.ShortId),0) as clickCount From urls 
			LEFT JOIN clicks
				ON urls.shortId = clicks.ShortId 
			WHERE urls.statId = ?
			GROUP BY urls.ShortID`
	row := d.db.QueryRow(query, statId)

	var shortID string
	var expirationDateStr string
	var clicksCount int64
	err = row.Scan(&shortID, &expirationDateStr, &clicksCount)
	if err == sql.ErrNoRows {
		error := errors.New("stat url doesn't exist")
		d.log.Error(error)
		return nil, error
	} else if err != nil {
		d.log.Error(err)
		return nil, err
	}

	expirationDate, _ := time.Parse("2006-01-02 15:04:05.999999999-07:00", expirationDateStr)

	query = `SELECT IP, Time FROM clicks WHERE ShortId = ? ORDER BY Time DESC LIMIT 100`
	rows, err := d.db.Query(query, shortID)

	var clicks []*models.ClickScheme
	defer rows.Close()
	for rows.Next() {
		var ip string
		var timeString string

		err := rows.Scan(&ip, &timeString)
		if err != nil {
			d.log.Error(err)
		}

		time, _ := time.Parse("2006-01-02 15:04:05.999999999-07:00", timeString)
		click := &models.ClickScheme{
			IP:   ip,
			Time: time.Format("2006-01-02 15:04:05"),
		}
		clicks = append(clicks, click)
	}

	ss = &models.StatsScheme{
		ClickCount:     clicksCount,
		ExpirationDate: expirationDate.Format("2006-01-02"),
		Clicks:         clicks,
	}

	return ss, nil
}

func getShortId(value int64) string {
	bi := big.NewInt(value)
	slice := bi.Bytes()
	return bytesToKey(&slice)
}

func NewStatKey() string {
	id, _ := uuid.NewUUID()
	slice := id[:]
	return bytesToKey(&slice)
}

func bytesToKey(src *[]byte) string {
	encodedString := b64.StdEncoding.EncodeToString((*src))
	encodedString = strings.TrimRight(encodedString, "=")
	return strings.Replace(encodedString, "/", "-", -1)
}
