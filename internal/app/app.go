package app

import (
	"log"
	"os"

	"github.com/augustoapg/censysKvStore/internal/api"
	"github.com/augustoapg/censysKvStore/internal/store"
)

type Application struct {
	Logger    *log.Logger
	KVHandler *api.KVHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	store := store.NewInMemoryKVStore()

	return &Application{
		Logger:    logger,
		KVHandler: api.NewKVHandler(store, logger),
	}, nil
}
