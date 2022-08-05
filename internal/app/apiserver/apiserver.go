package apiserver

import (
	"database/sql"
	"level_zero/internal/app/cache/lru_cache"
	"level_zero/store/sqlstore"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Start(config *Config) error {
	db, err := NewDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	logger := getNewLogger(config)
	stanInfo := config.GetNatsInfo()
	cache, err := lru_cache.NewCache(config.CacheSize, store, logger)
	if err != nil {
		log.Fatal("Failed to create cache")
	}

	srv := NewServer(store, cache, logger, &stanInfo)

	return http.ListenAndServe(config.BindAddr, srv)
}

func NewDB(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getNewLogger(config *Config) *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(config.LogLevel)
	if err == nil {
		logger.SetLevel(lvl)
	} else {
		logger.Error("logger: can't parse log level")
	}
	return logger
}
