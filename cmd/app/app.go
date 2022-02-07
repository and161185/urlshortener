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
	DBDriverName     string `yaml:"dbDriverName"`
	ConnectionString string `yaml:"connectionString"`
	LogLevel         string `yaml:"logLevel"`
	Port             int    `yaml:"port"`
	WriteTimeout     int    `yaml:"writetimeout"`
	ReadTimeout      int    `yaml:"readtimeout"`
}

type app struct {
	log    *logrus.Logger
	us     *usrepo.UrlShortener
	config *config
}

const defaultDBDriverName = "sqlite3"
const defaultConnectionString = "data.db"
const defaultLogLevel = logrus.InfoLevel
const defaultPort = 8080
const defaultWriteTimeout = 10
const defaultReadTimeout = 10

func getConfig(log *logrus.Logger, configPath string) *config {
	log.Info("loading settings")

	//для Heroku сделано получение конфигурации из переменных среды
	envPort, _ := strconv.Atoi(os.Getenv("PORT"))
	envWriteTimeout, _ := strconv.Atoi(os.Getenv("WRITETIMAOUT"))
	envReadTimeout, _ := strconv.Atoi(os.Getenv("READTIMEOUT"))
	cfg := &config{
		DBDriverName:     os.Getenv("DBDRIVERNAME"),
		ConnectionString: os.Getenv("CONNECTIONSTRING"),
		LogLevel:         os.Getenv("LOGLEVEL"),
		Port:             envPort,
		WriteTimeout:     envWriteTimeout,
		ReadTimeout:      envReadTimeout,
	}

	fileCfg, err := readConfigFile(log, configPath)

	if err != nil {

		log.Errorf("Couldn't read config file %s , got %v", configPath, err)
		log.Info("Use app's default settings")

		log.Level = defaultLogLevel
		log.Infof("logLevel's default value %v is setted", defaultLogLevel)

		fileCfg = &config{}
	} else {
		log.Infof("string logrus level: %s", cfg.LogLevel)
		level, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Errorf("Couldn't parse log level, got %v", err)
			level = defaultLogLevel
		}
		log.Level = level
	}

	if cfg.ConnectionString == "" {
		cfg.ConnectionString = fileCfg.ConnectionString
		if cfg.ConnectionString == "" {
			cfg.ConnectionString = defaultConnectionString
			log.Infof("DbName can't be empty. Default value %v is setted", defaultConnectionString)
		}
	}

	if cfg.DBDriverName == "" {
		cfg.DBDriverName = fileCfg.DBDriverName
		if cfg.DBDriverName == "" {
			cfg.DBDriverName = defaultDBDriverName
			log.Infof("DBDriverName can't be empty. Default value %v is setted", defaultDBDriverName)
		}
	}

	if cfg.Port == 0 {
		cfg.Port = fileCfg.Port
		if cfg.Port == 0 {
			cfg.Port = defaultPort
			log.Infof("Port can't be 0. Default value %v is setted", defaultPort)
		}
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = fileCfg.ReadTimeout
		if cfg.ReadTimeout == 0 {
			cfg.ReadTimeout = defaultReadTimeout
			log.Infof("ReadTimeout can't be 0. Default value %v is setted", defaultReadTimeout)
		}
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = fileCfg.WriteTimeout
		if cfg.WriteTimeout == 0 {
			cfg.WriteTimeout = defaultWriteTimeout
			log.Infof("WriteTimeout can't be 0. Default value %v is setted", defaultWriteTimeout)
		}
	}

	log.Info("Settings loaded")

	return cfg
}

func readConfigFile(log *logrus.Logger, configPath string) (*config, error) {

	log.Info("reading config file")

	result := &config{}

	_, err := os.Stat(configPath)

	if err != nil {
		log.Errorf("Couldn't stat config file %s , got %v", configPath, err)
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Couldn't read config file %s , got %v", configPath, err)
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, result)
	if err != nil {
		return nil, err
	}

	return result, nil
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

	conf := getConfig(log, (*configPath))

	a := &app{
		log:    log,
		config: conf,
	}

	log.Info("App initialized")

	return a
}

func (a *app) Run() {

	uss := usstorage.NewUSStorage(a.log, a.config.DBDriverName, a.config.ConnectionString)
	us := usrepo.NewUrlShortener(uss)
	defer uss.Close()

	a.us = us
	router := handler.NewHandler(a.log, us)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + strconv.Itoa(a.config.Port),
		WriteTimeout: time.Duration(a.config.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(a.config.ReadTimeout) * time.Second,
	}

	go func() {
		a.log.Infof("App is starting on port: %v", a.config.Port)
		a.log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.ReadTimeout+a.config.WriteTimeout)*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("err while shutting down", err)
	}
	a.log.Info("shutting down")
	os.Exit(0)
}
