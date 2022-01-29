package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"urlshortener/internal/api/handler"
	usstorage "urlshortener/internal/db"
	"urlshortener/internal/repos/usrepo"
)

type config struct {
	DbName       string `yaml:"dbname"`
	LogLevel     string `yaml:"logLevel"`
	Port         int    `yaml:"port"`
	WriteTimeout int    `yaml:"writetimeout"`
	ReadTimeout  int    `yaml:"readtimeout"`
}

type app struct {
	log    *logrus.Logger
	us     *usrepo.UrlShortener
	config *config
}

const defaultDbName = "data.db"
const defaultLogLevel = logrus.InfoLevel
const defaultPort = 80
const defaultWriteTimeout = 10
const defaultReadTimeout = 10

func readConfigFile(log *logrus.Logger, configPath string) *config {
	log.Info("loading settings")

	var logLevel logrus.Level
	result := &config{}

	_, err := os.Stat(configPath)

	defer func() {
		if err != nil {

			log.Errorf("Couldn't read config file %s , got %v", configPath, err)
			log.Info("Use app's default settings")

			logLevel = defaultLogLevel
			log.Infof("logLevel's default value %v is setted", defaultLogLevel)
		}

		if result.DbName == "" {
			result.DbName = defaultDbName
			log.Infof("Port can't be empty. Default value %v is setted", defaultDbName)
		}
		if result.Port == 0 {
			result.Port = defaultPort
			log.Infof("Port can't be 0. Default value %v is setted", defaultPort)
		}
		if result.ReadTimeout == 0 {
			result.ReadTimeout = defaultReadTimeout
			log.Infof("ReadTimeout can't be 0. Default value %v is setted", defaultReadTimeout)
		}
		if result.WriteTimeout == 0 {
			result.WriteTimeout = defaultWriteTimeout
			log.Infof("WriteTimeout can't be 0. Default value %v is setted", defaultWriteTimeout)
		}

		log.Level = logLevel
		log.Info("Settings loaded")

	}()

	if err != nil {
		return nil
	}

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(yamlFile, result)
	if err != nil {
		return nil
	}

	log.Infof("string logrus level: %s", result.LogLevel)
	level, err := logrus.ParseLevel(result.LogLevel)
	if err != nil {
		return nil
	}
	logLevel = level

	return result
}

func NewApp() *app {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	f, err := os.OpenFile(".\\log.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(logrus.TextFormatter)

	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {
		fmt.Println(err)
	} else {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}
	log.Formatter = new(logrus.JSONFormatter)

	var configPath *string = flag.String("conf", ".\\config\\config.yaml", "Configuration file's path")
	flag.Parse()

	conf := readConfigFile(log, (*configPath))

	a := &app{
		log:    log,
		config: conf,
	}

	log.Info("App initialized")

	return a
}

func (a *app) Run() {

	uss := usstorage.NewUSStorage(a.log, a.config.DbName)
	us := usrepo.NewUrlShortener(uss)
	defer us.Close()

	a.us = us
	router := handler.NewHandler(a.log, us)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + strconv.Itoa(a.config.Port),
		WriteTimeout: time.Duration(a.config.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(a.config.ReadTimeout) * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.ReadTimeout+a.config.WriteTimeout)*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	a.log.Info("shutting down")
	os.Exit(0)
}
