module urlshortener

// для Heroku определение версии go
// +heroku goVersion go1.17
go 1.17

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.10
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.2.2
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)
