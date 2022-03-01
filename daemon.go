package main

import (
	"belli/onki-game-ideas-mongo-backend/api"
	"belli/onki-game-ideas-mongo-backend/config"
	"belli/onki-game-ideas-mongo-backend/datastore"
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/service/cache"
	"belli/onki-game-ideas-mongo-backend/web"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	quitChan := make(chan interface{})
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "01-02 15:04:05.000"
	customFormatter.FullTimestamp = true

	log.SetFormatter(customFormatter)
	log.SetLevel(log.DebugLevel)

	err := config.ReadConfig("config.json")
	if err != nil {
		log.WithError(err).Panic("error reading config file")
	}

	err = errs.ReadErrorFile("errorprotocol.json")
	if err != nil {
		log.WithError(err).Panic("error reading error file")
	}

	setupServer(quitChan, signalChan, config.Conf)
}

func setupServer(quit chan interface{}, signalChan chan os.Signal, conf *config.Config) {
	access, err := datastore.NewMgAccess(conf)
	if err != nil {
		log.WithError(err).Panic("Could not initialize datastore.Access")
		return
	}

	cacheService := cache.NewRedisService(config.Conf.RedisConn, config.Conf.RedisDB, 100, 100)

	apiController := api.NewAPIController(access, cacheService)

	s := web.NewServer(apiController)

	r := web.NewRouter(s)

	srv := &http.Server{
		Addr:         conf.ListenAddress,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 45 * time.Second,
	}

	listener, err := net.Listen("tcp", conf.ListenAddress)
	if err != nil {
		log.WithError(err).Panic(" setting up listener")
		return
	}

	log.WithField("listen", conf.ListenAddress).Info("Starting HTTP API Server")

	fmt.Println("<--START-SERVER-->")

	go startServer(srv, listener)

	for {
		select {
		case <-quit:
			log.Warn("quit channel closed, closing listener")
			err = srv.Close()
			if err != nil {
				log.WithError(err).Error("error during HTTP Server close")
			}
			err = listener.Close()
			if err != nil {
				log.WithError(err).Error("error during TCP Listener close")
			}
			return
		case sig := <-signalChan:
			switch sig {
			case os.Interrupt, os.Kill, syscall.SIGTERM:
				log.Info("interrupt signal received, sending Quit signal")
				close(quit)
			default:
				log.WithField("signal", sig).Info("signal received")
			}
		}
	}
}

func startServer(srv *http.Server, listener net.Listener) {
	err := srv.Serve(listener)
	if err != nil {
		log.WithError(err).Error("HTTP server Error")
	}
}
