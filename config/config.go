package config

import (
	"belli/onki-game-ideas-mongo-backend/helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type IdeaFile struct {
	MaxSize           int64    `json:"max_size"`
	AllowedExtensions []string `json:"allowed_extensions"`
	ContentTypes      []string `json:"content_types"`
}

type BlockExpiry struct {
	IP   int `json:"ip"`
	User int `json:"user"`
}
type Jwt struct {
	AccessTokenExpiry  int    `json:"access_token_expiry"`
	RefreshTokenExpiry int    `json:"refresh_token_expiry"`
	Secret             string `json:"secret"`
}

type PrefixForRedis struct {
	RefreshToken             string `json:"refresh_token"`
	WorkerRestrictIdeaSubmit string `json:"worker_resrtict"`
}

type Config struct {
	DbConn         string         `json:"db_conn"`
	MgConn         string         `json:"mg_conn"`
	RedisConn      string         `json:"redis_conn"`
	RedisDB        int            `json:"redis_db"`
	ListenAddress  string         `json:"listen_address"`
	IdeaFile       *IdeaFile      `json:"idea_file"`
	Jwt            *Jwt           `json:"jwt"`
	Static         string         `json:"static"`
	StaticDir      string         `json:"static_dir"`
	Renderer       string         `json:"renderer_url"`
	PrefixForRedis PrefixForRedis `json:"prefix_for_redis"`
}

type ConfigFile struct {
	IdeaFile       *IdeaFile      `json:"idea_file"`
	Jwt            *Jwt           `json:"jwt"`
	Static         string         `json:"static"`
	StaticDir      string         `json:"static_dir"`
	PrefixForRedis PrefixForRedis `json:"prefix_for_redis"`
}

var Conf *Config

func ReadConfig(source string) (err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}

	var raw []byte
	raw, err = ioutil.ReadFile(source)
	if err != nil {
		wMsg := "error reading config from file, creating new sample"
		log.Warn(wMsg)

		err = createDefaultConfig(source)
		if err != nil {
			eMsg := "error creating config file"
			log.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		raw, err = ioutil.ReadFile(source)
		if err != nil {
			eMsg := "error reading config from file"
			log.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
	}
	confFile := &ConfigFile{}
	err = json.Unmarshal(raw, &confFile)
	if err != nil {
		eMsg := "error parsing config from json"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		Conf = nil
		return
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		eMsg := "error on converting REDIS_DB to int"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	Conf = &Config{
		DbConn: fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
		),
		MgConn:         os.Getenv("MG_CONN_LOCAL"),
		Renderer:       os.Getenv("RENDERER_URL"),
		RedisConn:      fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		RedisDB:        redisDB,
		ListenAddress:  os.Getenv("LISTEN_ADDR"),
		IdeaFile:       confFile.IdeaFile,
		Jwt:            confFile.Jwt,
		Static:         confFile.Static,
		StaticDir:      confFile.StaticDir,
		PrefixForRedis: confFile.PrefixForRedis,
	}

	return
}

func createDefaultConfig(source string) (err error) {
	static := filepath.Join(filepath.Dir(source), "static")

	secret, err := helpers.GenRandomString(64)
	if err != nil {
		eMsg := "error on generate random secret"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	c := ConfigFile{
		IdeaFile: &IdeaFile{
			MaxSize: 10485760,
			AllowedExtensions: []string{
				"jpg",
				"png",
				"jpeg",
			},
			ContentTypes: []string{
				"image/jpeg",
				"image/png",
			},
		},
		PrefixForRedis: PrefixForRedis{
			RefreshToken: "USER_REFRESH_TOKEN:",
		},

		Jwt: &Jwt{
			AccessTokenExpiry:  3600,
			RefreshTokenExpiry: 86400,
			Secret:             secret,
		},

		Static:    static,
		StaticDir: static,
	}

	b, err := json.MarshalIndent(c, "", "    ")

	if err != nil {
		eMsg := "error marshall config file"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	err = ioutil.WriteFile(source, b, 0644)
	if err != nil {
		eMsg := "error creating config file"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}
