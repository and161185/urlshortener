package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"urlshortener/internal/models"
	"urlshortener/internal/repos/usrepo"
)

func NewHandler(log *logrus.Logger, repo *usrepo.UrlShortener) *mux.Router {
	router := mux.NewRouter()

	handlerGenerate := &HandlerGenerate{log: log, repo: repo}
	router.HandleFunc("/generate", handlerGenerate.ServeHTTP).Methods("POST")

	handlerStats := &HandlerStats{log: log, repo: repo}
	router.HandleFunc("/stat/{statid}", handlerStats.ServeHTTP).Methods("GET")

	handlerRedirect := &HandlerRedirect{log: log, repo: repo}
	router.HandleFunc("/{shorturl}", handlerRedirect.ServeHTTP).Methods("GET")

	loggingMiddleware := LoggingMiddleware(log)
	router.Use(loggingMiddleware)

	return router
}

type HandlerGenerate struct {
	log  *logrus.Logger
	repo *usrepo.UrlShortener
}

func (h *HandlerGenerate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HandlerGenerate")

	var urlData models.FullUrlScheme
	err := json.NewDecoder(r.Body).Decode(&urlData)
	if err != nil {
		h.log.Error(err)
		fmt.Fprint(w, err)
		return
	}

	data, err := h.repo.GenerateShortUrl(urlData)
	if err != nil {
		h.log.Error(err)
		fmt.Fprint(w, err)
		return
	}

	h.log.Debug(data)

	bytes, err := json.Marshal(data)
	if err != nil {
		h.log.Error(err)
		fmt.Fprint(w, err)
		return
	}
	answer := string(bytes)

	h.log.Debug(answer)
	fmt.Fprint(w, answer)
}

type HandlerStats struct {
	log  *logrus.Logger
	repo *usrepo.UrlShortener
}

func (h *HandlerStats) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HandlerStats")

	statId := strings.TrimLeft(r.RequestURI, "/stat/")

	statsStruct, err := h.repo.GetStats(statId)
	if err != nil {
		h.log.Error(err)
		fmt.Fprint(w, err)
		return
	}

	h.log.Debug(statsStruct)

	bytes, err := json.Marshal(statsStruct)
	if err != nil {
		h.log.Error(err)
		fmt.Fprint(w, err)
		return
	}
	answer := string(bytes)

	h.log.Debug(answer)
	fmt.Fprint(w, answer)
}

type HandlerRedirect struct {
	log  *logrus.Logger
	repo *usrepo.UrlShortener
}

func (h *HandlerRedirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	shortId := strings.TrimLeft(r.RequestURI, "/")

	urlScheme, err := h.repo.GetFullUrl(shortId)
	if err != nil {
		h.log.Error(err)
		fmt.Fprint(w, err)
		return
	}
	url := urlScheme.Url

	http.Redirect(w, r, url, http.StatusMovedPermanently)
	go h.registerClick(shortId, r)
}

func (h *HandlerRedirect) registerClick(shortId string, r *http.Request) {

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.log.Errorf("can't split %s", r.RemoteAddr)
		ip = "undefined"
	} else {
		netip := net.ParseIP(ip)
		if netip == nil {
			h.log.Errorf("can't get ip from %s", r.RemoteAddr)
			ip = "undefined"
		}
	}

	h.repo.RegisterClick(shortId, ip)
}

func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Errorln(err)
				}
			}()

			logger.Infoln("method", r.Method, "path", r.URL.EscapedPath())
			next.ServeHTTP(w, r)

		}

		return http.HandlerFunc(fn)
	}
}
