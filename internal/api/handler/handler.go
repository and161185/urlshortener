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

	"github.com/rs/cors"
)

func NewHandler(log *logrus.Logger, repo *usrepo.UrlShortener) http.Handler {
	router := mux.NewRouter()

	handler := &Handler{log: log, repo: repo}
	router.HandleFunc("/generate", handler.generate).Methods("POST")

	router.HandleFunc("/stat/{statid}", handler.stat).Methods("GET")

	router.HandleFunc("/{shorturl}", handler.redirect).Methods("GET")

	router.HandleFunc("/heart/beat", handler.heartbeat).Methods("GET")

	router.HandleFunc("/", handler.front).Methods("GET")

	CorsHandler := cors.Default().Handler(router)
	loggingMiddleware := LoggingMiddleware(log)

	router.Use(loggingMiddleware)

	return CorsHandler
}

type Handler struct {
	log  *logrus.Logger
	repo *usrepo.UrlShortener
}

func (h *Handler) heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (h *Handler) generate(w http.ResponseWriter, r *http.Request) {

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

func (h *Handler) stat(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) registerClick(shortId string, r *http.Request) {

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

	err = h.repo.RegisterClick(shortId, ip)
	if err != nil {
		h.log.Errorf("can't register click %s from ip %s", shortId, ip)
	}
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
